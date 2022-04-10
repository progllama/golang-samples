package file_downloader

import (
	"io"
	"net/http"
	"os"
)

func FileDownloader() {
	downloadFile("./samples/file_downloader/files/test.mp4", "https://cc3001.dmm.co.jp/litevideo/freepv/h/h_0/h_086cherd00082/h_086cherd00082_dmb_w.mp4")
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
