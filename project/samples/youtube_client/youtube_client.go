package youtube_client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var SEARCH_KEYWORDS []string = []string{
	"cute",
	"nice",
	"good",
	"beautiful",
	"anime",
	"furry",
	"guide",
	"tutorial",
	"modeling",
	"simple",
	"begginer",
	"plain",
	"easy",
	"lopoly",
	"female",
	"girl",
	"woman",
	"human",
	"character",
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func RunYoutubeClient() {
	loadEnv()

	key := os.Getenv("YOUTUBE_API_KEY")

	ctx := context.Background()
	client, err := youtube.NewService(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Fatal("Error creating new Youtube client: %v", err)
	}

	// service, err := youtube.New
	// if err != nil {
	// 	log.Fatal("Error creating new Youtube client: %v", err)
	// }

	call := service.Search.List([]string{"id, snippet"}).
		Q("Blender").
		PublishedAfter("2022-04-07T00:00:00Z").
		PublishedBefore("2022-04-08T00:00:00Z").
		MaxResults(100)

	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.ChannelId
		case "youtube#playlists":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	nextPageToken := response.NextPageToken

	call = service.Search.List([]string{"id, snippet"}).
		Q("Blender").
		PublishedAfter("2022-04-07T00:00:00Z").
		PublishedBefore("2022-04-08T00:00:00Z").
		PageToken(nextPageToken).
		MaxResults(100)

	response, err = call.Do()
	if err != nil {
		log.Fatal(err)
	}

	printIDs("Videos", videos)
	printIDs("Channels", channels)
	printIDs("Playlists", playlists)
}

func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, _ := range matches {
		fmt.Printf("https://www.youtube.com/watch?v=%v\n\n", id)
	}
	fmt.Printf("\n\n")
}
