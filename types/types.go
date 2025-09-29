package types

import "encoding/xml"

type ListBucketResult struct {
	XMLName     xml.Name   `xml:"ListBucketResult"`
	Contents    []Contents `xml:"Contents"`
	Prefix      string     `xml:"Prefix"`
	IsTruncated bool       `xml:"IsTruncated"`
}

type Contents struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int64  `xml:"Size"`
	StorageClass string `xml:"StorageClass"`
}

func (c *Contents) ObjURL() string {
	return baseUrl + "/" + c.Key
}
