package mspyls

import (
	"encoding/xml"
	"sort"
	"time"

	"github.com/johejo/extratime"
)

type EnumerationResults struct {
	XMLName       xml.Name `xml:"EnumerationResults"`
	Text          string   `xml:",chardata"`
	ContainerName string   `xml:"ContainerName,attr"`
	Blobs         Blobs    `xml:"Blobs"`
	NextMarker    string   `xml:"NextMarker"`
}

type Blobs struct {
	Text string `xml:",chardata"`
	Blob []Blob `xml:"Blob"`
}

type Blob struct {
	Text       string     `xml:",chardata"`
	Name       string     `xml:"Name"`
	URL        string     `xml:"Url"`
	Properties Properties `xml:"Properties"`
}

type Properties struct {
	Text            string            `xml:",chardata"`
	LastModified    extratime.RFC1123 `xml:"Last-Modified"`
	Etag            string            `xml:"Etag"`
	ContentLength   uint64            `xml:"Content-Length"`
	ContentType     string            `xml:"Content-Type"`
	ContentEncoding string            `xml:"Content-Encoding"`
	ContentLanguage string            `xml:"Content-Language"`
	ContentMD5      string            `xml:"Content-MD5"`
	CacheControl    string            `xml:"Cache-Control"`
	BlobType        string            `xml:"BlobType"`
	LeaseStatus     string            `xml:"LeaseStatus"`
}

var (
	_ sort.Interface = (*Blobs)(nil)
)

func (b *Blobs) Len() int {
	return len(b.Blob)
}

func (b *Blobs) Less(i, j int) bool {
	ti := time.Time(b.Blob[i].Properties.LastModified).UnixNano()
	tj := time.Time(b.Blob[j].Properties.LastModified).UnixNano()
	return ti < tj
}

func (b *Blobs) Swap(i, j int) {
	b.Blob[i], b.Blob[j] = b.Blob[j], b.Blob[i]
}
