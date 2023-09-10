package utils

import (
	"github.com/qianlnk/pgbar"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var pgb = pgbar.New("")

func NewFile(url string) *File {
	return &File{
		started: time.Now(),
		Url:     url,
	}
}

type Downloader struct {
	io.Reader
	bar *pgbar.Bar
}
type File struct {
	Url              string
	started          time.Time
	Total            int64
	Done             int
	DownloadComplete bool
}

func (f *File) GetPath() string {
	parser, _ := url.Parse(f.Url)
	return filepath.Join(parser.Hostname(), parser.Path)
}

func (f *File) CreatePath() {
	err := os.MkdirAll(filepath.Dir(f.GetPath()), 0755)
	if err != nil {
		return
	}
}

func (f *File) StartDownload() error {

	f.CreatePath()

	filePath := f.GetPath()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(f.Url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bar := pgbar.NewBar(0, "[D] Downloading :", int(resp.ContentLength))
	if resp.ContentLength > 10*1024 {
		bar.SetUnit("B", "kb", 1024*1024)
	} else if resp.ContentLength > 10*1024*1024 {
		bar.SetUnit("B", "MB", 1024*1024)
	}
	_, err = io.Copy(file, &Downloader{
		Reader: resp.Body,
		bar:    bar,
	})
	if err != nil {
		return err
	}
	f.DownloadComplete = true
	return nil
}

func (d *Downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)
	d.bar.Add(n)
	return n, err
}
