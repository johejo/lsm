package mspyls

import (
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnumeratedResult(t *testing.T) {
	f, err := os.Open("./testdata/test.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var er EnumerationResults
	if err := xml.NewDecoder(f).Decode(&er); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, time.Time(er.Blobs.Blob[0].Properties.LastModified).Format(time.RFC1123), "Tue, 17 Dec 2019 17:41:17 GMT")
}
