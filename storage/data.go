package storage

import (
	"github.com/MegaBytee/binance-data-go/types"
	"gorm.io/gorm/clause"
)

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
