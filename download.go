package binance_data

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (d *Downloader) downloadFile(ctx context.Context, u string) (string, error) {
	local, err := localPathForURL(u)
	if err != nil {
		return "", err
	}
	// ensure dir
	if err := ensureDir(filepath.Dir(local)); err != nil {
		return "", err
	}

	// HEAD to check size
	req, _ := http.NewRequestWithContext(ctx, "HEAD", u, nil)
	if resp, err := d.httpClient.Do(req); err == nil && resp != nil {
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if sizeS := resp.Header.Get("Content-Length"); sizeS != "" {
				if fi, err := os.Stat(local); err == nil {
					if fi.Size() == parseInt64(sizeS) {
						return local, fmt.Errorf("exists")
					}
				}
			}
		}
	}

	// GET
	req2, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	resp2, err := d.httpClient.Do(req2)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		return "", fmt.Errorf("bad status %d", resp2.StatusCode)
	}

	tmp := local + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return "", err
	}
	n, err := io.Copy(out, resp2.Body)
	out.Close()
	if err != nil {
		os.Remove(tmp)
		return "", err
	}
	if err := os.Rename(tmp, local); err != nil {
		return "", err
	}
	fmt.Printf("Downloaded: %s (%d bytes)\n", local, n)
	return local, nil
}
