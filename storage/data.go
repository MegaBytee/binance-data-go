package storage

import (
	"github.com/MegaBytee/binance-data-go/types"
	"gorm.io/gorm/clause"
)

/*
import (

	"gorm.io/gorm/clause"

)

	func (s *Storage) SaveSite(site *types.Site) error {
		s.mu.Lock()
		defer s.mu.Unlock()

		return s.Data.Clauses(clause.OnConflict{DoNothing: true}).Create(site).Error
	}

	func (s *Storage) CreateSitemapsInBatches(sitemap []types.SitemapIndex) error {
		s.mu.Lock()
		defer s.mu.Unlock()

		return s.Data.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&sitemap, 100).Error

}

	func (s *Storage) GetSitemapIndexToScan(limit int) []types.SitemapIndex {
		s.mu.RLock()
		defer s.mu.RUnlock()

		var sitemaps []types.SitemapIndex
		s.Data.Limit(limit).Where("scanned=?", false).Find(&sitemaps)
		return sitemaps
	}

	func (s *Storage) UpdateScannedSitemaps(sitemaps []types.SitemapIndex) error {
		s.mu.Lock()
		defer s.mu.Unlock()
		return updateOneFieldInBatches(s.Data,
			&types.SitemapIndex{}, "hash IN ?",
			types.GetHashIDs(sitemaps), "scanned", true, 100)
	}

	func (s *Storage) CreateLinksInBatches(links []types.Link) error {
		s.mu.Lock()
		defer s.mu.Unlock()

		return s.Data.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&links, 100).Error

}
*/
func (s *Storage) CreateFilesInBatches(files []types.File) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.Data.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&files, 100).Error

}
func (s *Storage) GetFilesByStatus(status types.FileStatus, limit int) []types.File {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var files []types.File
	s.Data.Limit(limit).Where("status=?", int(status)).Find(&files)
	return files
}

func (s *Storage) UpdateFile(file types.File) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.Data.Model(&file).Where("hash=?", file.Hash).
		Select("status", "local").
		Updates(types.File{Status: file.Status, Local: file.Local}).Error

}
func (s *Storage) UpdateExtractedFiles(files []types.File) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return updateOneFieldInBatches(s.Data,
		&types.File{}, "hash IN ?",
		types.GetHashIDs(files), "status", int(types.FileStatusExtracted), 100)
}
