package file_downloader

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func FileDownloader() {
	in := "./samples/web_video_player/web_video_player_save_file/video_urls.txt"
	paths := getPaths(in)
	for i, p := range paths {
		out := fmt.Sprintf("./samples/file_downloader/files/%v.mp4", i)
		downloadFile(out, p)
	}
}

func getPaths(in string) []string {
	f, err := os.Open(in)
	if err != nil {
		return make([]string, 0)
	}

	paths := make([]string, 0)
	bs := bufio.NewScanner(f)
	for bs.Scan() {
		paths = append(paths, bs.Text())
	}
	return paths
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
