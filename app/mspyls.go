package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/johejo/extratime"
)

type MSPyLSInstaller struct {
	baseInstaller
}

var _ Installer = (*MSPyLSInstaller)(nil)

func NewMSPyLSInstaller(baseDir string) *MSPyLSInstaller {
	var i MSPyLSInstaller
	i.baseInstaller = newBaseInstaller(filepath.Join(baseDir, i.Name()))
	return &i
}

func (i *MSPyLSInstaller) Name() string {
	return "microsoft-python-language-server"
}

func (i *MSPyLSInstaller) BinName() string {
	return noExecutable
}

func (i *MSPyLSInstaller) Version() string {
	return versionUnSpecified
}

func (i *MSPyLSInstaller) Requires() []string {
	return noRequires
}

const xmlURL = "https://pvsc.blob.core.windows.net/python-language-server-stable?restype=container&comp=list&prefix=Python-Language-Server"

func (i *MSPyLSInstaller) Install(ctx context.Context) error {
	if err := i.downloadNuPkg(ctx); err != nil {
		return err
	}
	if err := i.installDotNet(ctx); err != nil {
		return err
	}
	return nil
}

func (i *MSPyLSInstaller) downloadNuPkg(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, xmlURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("invalid status %v, err=%v", resp.StatusCode, err)
		}
		return fmt.Errorf("invalid status %v, body=%v", resp.StatusCode, string(b))
	}

	var r EnumerationResults
	if err := xml.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	blobs := r.Blobs.FilterByOS()
	sort.Sort(&blobs)
	latest := blobs.Blob[len(blobs.Blob)-1]
	u := latest.URL
	log.Printf("download nupkg from %s", u)
	if err := i.FetchWithExtract(ctx, u, filepath.Join(i.Dir(), i.Name()+".zip")); err != nil {
		return err
	}
	return nil
}

func (i *MSPyLSInstaller) installDotNet(ctx context.Context) error {
	var (
		dotNetInstallerURL, _script string
		args                        []string
	)
	switch runtime.GOOS {
	case linux, darwin:
		dotNetInstallerURL = "https://dot.net/v1/dotnet-install.sh"
		_script = "dotnet-install.sh"
		args = []string{"--install-dir", filepath.Join(i.Dir(), "dotnet")}
	case windows:
		dotNetInstallerURL = "https://dot.net/v1/dotnet-install.ps1"
		_script = "dotnet-install.ps1"
		args = []string{"-InstallDir", filepath.Join(i.Dir(), "dotnet")}
	default:
		return fmt.Errorf("%s is not supported dotnet install scripts", runtime.GOOS)
	}

	log.Printf("download %s from %s", _script, dotNetInstallerURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dotNetInstallerURL, nil)
	if err != nil {
		return err
	}
	script := filepath.Join(i.Dir(), _script)
	if err := i.Download(req, script); err != nil {
		return err
	}
	if err := os.Chmod(script, 0755); err != nil {
		return err
	}
	if err := i.CmdRun(ctx, script, args...); err != nil {
		return err
	}
	return nil
}

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

func (b *Blobs) FilterByOS() Blobs {
	var _os string
	switch runtime.GOOS {
	case windows:
		_os = "win"
	case darwin:
		_os = "osx"
	case linux:
		_os = "linux"
	default:
		_os = "any"
	}
	filtered := make([]Blob, 0, len(b.Blob))
	for _, blob := range b.Blob {
		if strings.Contains(blob.URL, _os) {
			filtered = append(filtered, blob)
		}
	}
	return Blobs{Text: b.Text, Blob: filtered}
}
