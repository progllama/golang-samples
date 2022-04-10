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
			// "--headless",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)
	err := driver.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Stop()
	page, err := driver.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	// 自動操作
	urls := getURL()
	for _, url := range urls {
		fmt.Println(url)
		err = page.Navigate(url)
		if err != nil {
			log.Fatal(err)
			return
		}

		s := page.FindByLink("はい")
		s.Click()
		// s = page.FindByLink("無料サンプル動画を見る")
		// var result string
		// page.RunScript("document.querySelectorAll('#detail-sample-movie a')[0].click();", nil, &result)
		// err = s.DoubleClick()
		// time.Sleep(10)
		s = page.Find("iframe:nth-child(1)")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		s.SwitchToFrame()
		s = page.FindByButton("再生/一時停止(Space)")
		s.Click()
		s.Click()
		s = page.Find(".player")
		text, _ := s.Text()
		t := strings.Split(text, "\n")[1]
		t = strings.Split(t, "/")[1]
		t = strings.TrimSpace(t)
		hour, _ := strconv.Atoi(strings.Split(t, ":")[0])
		hour = hour * 3600
		min, _ := strconv.Atoi(strings.Split(t, ":")[1])
		min = min * 60
		sec, _ := strconv.Atoi(strings.Split(t, ":")[2])
		time.Sleep(time.Duration((hour + min + sec)) * time.Second)
	}
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
