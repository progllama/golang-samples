package main

import (
	"fmt"
	"gorm_sample/samples/youtube_client"
)

func main() {
	fmt.Println("This is sample.")
	youtube_client.RunYoutubeClient()
	// goquery_sample.GoquerySample()
	// gorm_sample.GormMain()
	// collectpostalcode.CollectCountryNamesByGoquery()
}
