package lazy_image

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"io"
	"os"
	"path"
	"r34-client/commons"
	"sync"
)

var l = commons.NewLogger("lazy_img.downloader: ")

var downloader struct {
	downloadPath string
	// todo: perfile lock? how?
	fileLocks sync.Mutex
}

func init() {
	// initialize downloader path
	usrCacheDir, err := os.UserCacheDir()
	if err != nil {
		l.Fatalf("error getting user cache dir: %s", err)
	}
	downloader.downloadPath = path.Join(usrCacheDir, "r34-client")
	l.Printf("will use this path for cache: %s", downloader.downloadPath)
	createCacheDirectory(downloader.downloadPath)
	// initialize locks
	downloader.fileLocks = sync.Mutex{}
}

func createCacheDirectory(cachePath string) {
	err := os.MkdirAll(cachePath, os.ModePerm)
	if err != nil {
		l.Fatalf("error MkdirAll %s: %s", cachePath, err)
	}
}

func DownloadAndGetOutputFilePath(url string) (string, error) {
	// acquire lock
	downloader.fileLocks.Lock()
	defer downloader.fileLocks.Unlock()

	uri, err := storage.ParseURI(url)
	if err != nil {
		return "", fmt.Errorf("parse image url: %s, err: %s", url, err)
	}

	cachedLocalFileURI := genHashFileURI(uri)
	l.Printf("url to cached file url map: %s => %s", uri, cachedLocalFileURI)

	// check existing file is exists
	exist, err := storage.Exists(cachedLocalFileURI)
	if err != nil {
		return "", fmt.Errorf("error checking cached local file %s exists: %s", cachedLocalFileURI, err)
	}
	if exist {
		return cachedLocalFileURI.Path(), nil
	}

	reader, err := storage.Reader(uri)
	if err != nil {
		return "", fmt.Errorf("get image url: %s, err: %s", url, err)
	}
	defer reader.Close()

	writer, err := storage.Writer(cachedLocalFileURI)
	if err != nil {
		return "", fmt.Errorf("error create writer for: %s, err: %s", cachedLocalFileURI, err)
	}
	defer writer.Close()
	_, err = io.Copy(writer, reader)
	if err != nil {
		return "", fmt.Errorf("error doing io.Copy for: %s, err: %s", cachedLocalFileURI, err)
	}

	return cachedLocalFileURI.Path(), nil
}

func genHashFileURI(uri fyne.URI) fyne.URI {
	hashed := sha256.Sum256([]byte(uri.String()))
	hashedStr := hex.EncodeToString(hashed[:])
	outFileName := fmt.Sprintf("%s%s", hashedStr, uri.Extension())

	return storage.NewFileURI(path.Join(downloader.downloadPath, outFileName))
}

func ClearCache() error {
	downloader.fileLocks.Lock()
	defer downloader.fileLocks.Unlock()
	defer createCacheDirectory(downloader.downloadPath)
	return os.RemoveAll(downloader.downloadPath)
}
