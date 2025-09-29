package binance_data

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/MegaBytee/binance-data-go/config"
	"github.com/MegaBytee/binance-data-go/storage"
	"github.com/MegaBytee/binance-data-go/types"
	"github.com/gocolly/colly/v2"
)

const (
	httpTimeout = 30 * time.Second
)

type Downloader struct {
	c          *colly.Collector
	db         *storage.Storage
	httpClient *http.Client
}

func NewDownloader(cfg *config.Config) *Downloader {
	c, err := newCollyScrapper(cfg)
	if err != nil {
		return nil
	}
	db := storage.New().Config()
	if db == nil {
		return nil
	}
	return &Downloader{
		c:          c,
		db:         db,
		httpClient: &http.Client{Timeout: httpTimeout},
	}
}
func (d *Downloader) Close() {
	d.db.Close()
}

func (d *Downloader) Run(params *types.DataParams) {
	start := time.Now()
	ctx := context.Background()

	files := d.getFiles(params)
	if err := d.db.CreateFilesInBatches(files); err != nil {
		log.Printf("something went wrong:in  %v", err)
	}

	d.downloadFiles(ctx, 10)

	d.extractFiles(100)

	elapsed := time.Since(start)

	log.Println("Downloader executed in :", elapsed)
}

func (d *Downloader) getFiles(params *types.DataParams) []types.File {
	url := params.GetBucketUrl()
	d.c.AllowURLRevisit = false
	var body []byte
	d.c.OnResponse(func(r *colly.Response) {
		body = r.Body
	})

	err := d.c.Visit(url)
	if err != nil {
		log.Printf("something went wrong:err %v", err)
		return nil
	}
	var result types.ListBucketResult
	if err := xml.Unmarshal(body, &result); err != nil {
		log.Fatalf("xml unmarshal: %v", err)
	}

	return types.NewFiles(result.Contents, params)
}

func (d *Downloader) downloadFiles(ctx context.Context, limit int) {

	for {
		time.Sleep(time.Duration(3) * time.Second)
		files := d.db.GetFilesByStatus(types.FileStatusNew, limit)
		if len(files) > 0 {
			fmt.Println("downloading files....")

			for _, file := range files {
				delay := rand.Intn(9) + 3 // Random number between x and y
				fmt.Printf("Sleeping for %d s before the next request...\n", delay)
				local, err := d.downloadFile(ctx, file.Link)
				if err != nil {
					if err.Error() == "exists" {
						fmt.Printf("Skipping exists: %s\n", file.Link)
						continue
					}
					fmt.Printf("Download error %s: %v\n", file.Link, err)
					continue
				}
				//update db file
				file.Status = int(types.FileStatusDownloaded)
				file.Local = local
				d.db.UpdateFile(file)
				time.Sleep(time.Duration(delay) * time.Second)
			}

		} else {
			break
		}

	}

}

func (d *Downloader) extractFiles(limit int) {

	for {
		time.Sleep(time.Duration(5) * time.Second)
		files := d.db.GetFilesByStatus(types.FileStatusDownloaded, limit)
		if len(files) > 0 {
			fmt.Println("extracting files....")
			// Create a WaitGroup to wait for all goroutines to finish
			var wg sync.WaitGroup
			// Create a buffered channel to limit concurrency
			sem := make(chan struct{}, 100) // Limit concurrent goroutines

			for _, file := range files {
				wg.Add(1) // Increment the WaitGroup counter

				go func(file types.File) {
					defer wg.Done()   // Notify that this goroutine is done
					sem <- struct{}{} // Acquire a token
					if err := extractZipFile(file.Local); err != nil {
						fmt.Printf("Extract error %s: %v\n", file.Link, err)
					}
					<-sem // Release the token
				}(file)
			}

			// Wait for all goroutines to finish
			wg.Wait()
			d.db.UpdateExtractedFiles(files)
		} else {
			break
		}

	}
}
