package types

import "strings"

type FileStatus int

const (
	FileStatusNew FileStatus = iota
	FileStatusDownloaded
	FileStatusExtracted
	FileStatusError
)

type File struct {
	Hash   string `json:"hash" gorm:"unique"`
	From   string `json:"from"`
	Key    string `json:"key"`
	Size   int64  `json:"size"`
	Link   string `json:"link"`
	Status int    `json:"status"`
	Local  string
}

func (f File) HashID() string {
	return f.Hash
}
func NewFile(item Contents, params *DataParams) *File {
	return &File{
		Hash:   Hash256(item.ObjURL()),
		From:   params.From,
		Link:   item.ObjURL(),
		Key:    item.Key,
		Size:   item.Size,
		Status: int(FileStatusNew),
	}
}

func NewFiles(items []Contents, params *DataParams) []File {
	var files []File
	for _, item := range items {
		if strings.HasSuffix(item.Key, ".CHECKSUM") {
			continue
		}
		newFile := NewFile(item, params)
		files = append(files, *newFile)

	}
	return files
}
