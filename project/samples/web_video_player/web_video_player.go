package web_video_player

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

func WebVideoPlayer() {
	// Chromeを利用することを宣言
	var driver *agouti.WebDriver = agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)
	err := driver.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Stop()

	// page, err := driver.NewPage()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// 自動操作
	urls := getURL()
	data := make([]string, 0)
	for _, url := range urls {

		page, err := driver.NewPage()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)
		err = page.Navigate(url)
		if err != nil {
			log.Fatal(err)
			return
		}

		s := page.FindByLink("はい")
		s.Click()
		s = page.Find("iframe:nth-child(1)")
		s.SwitchToFrame()
		// s = page.FindByButton("再生/一時停止(Space)")
		// s.Click()
		// s.Click()

		// s = page.Find(".player")
		// text, _ := s.Text()
		// duration := parseTime(text)
		// time.Sleep(duration)

		s = page.Find("video")
		txt, _ := s.Attribute("src")
		data = append(data, txt)
		page.CloseWindow()
	}
	saveText(data)
}

func saveText(text []string) {
	path := "./samples/web_video_player/web_video_player_save_file/video_urls.txt"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	for _, line := range text {
		bw.WriteString(line + "\n")
	}
	bw.Flush()
}

func parseTime(txt string) time.Duration {
	t := strings.Split(txt, "\n")[1]
	t = strings.Split(t, "/")[1]
	t = strings.TrimSpace(t)

	hour, _ := strconv.Atoi(strings.Split(t, ":")[0])
	hour = hour * 3600
	min, _ := strconv.Atoi(strings.Split(t, ":")[1])
	min = min * 60
	sec, _ := strconv.Atoi(strings.Split(t, ":")[2])

	return time.Duration(hour+min+sec) * time.Second
}

func getURL() []string {
	f, err := os.Open("./samples/web_video_player/web_video_player.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	urls := make([]string, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		urls = append(urls, s.Text())
	}
	return urls
}

const DOWNLOAD_JS = `
var canvas = document.querySelector("canvas");

// Optional frames per second argument.
var stream = canvas.captureStream(25);
var recordedChunks = [];

console.log(stream);
var options = { mimeType: "video/webm; codecs=vp9" };
mediaRecorder = new MediaRecorder(stream, options);

mediaRecorder.ondataavailable = handleDataAvailable;
mediaRecorder.start();

function handleDataAvailable(event) {
  console.log("data-available");
  if (event.data.size > 0) {
    recordedChunks.push(event.data);
    console.log(recordedChunks);
    download();
  } else {
    // ...
  }
}
function download() {
  var blob = new Blob(recordedChunks, {
    type: "video/webm"
  });
  var url = URL.createObjectURL(blob);
  var a = document.createElement("a");
  document.body.appendChild(a);
  a.style = "display: none";
  a.href = url;
  a.download = "test.webm";
  a.click();
  window.URL.revokeObjectURL(url);
}

// demo: to download after 9sec
setTimeout(event => {
  console.log("stopping");
  mediaRecorder.stop();
}, 9000);
`
