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

const YOUTUBE_API_KEY = "YOUTUBE_API_KEY_VAR_NAME"

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getApiKey() string {
	return os.Getenv(YOUTUBE_API_KEY)
}

func RunYoutubeClient() {
	loadEnv()

	key := getApiKey()

	callYoutubeAPI(key, "blender", "2022-04-07T00:00:00Z", "2022-04-08T00:00:00Z", 100)
}

func makeService(key string) *youtube.Service {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Fatal("Error creating new Youtube client: %v", err)
	}
	return service
}

func makeApiCall(service *youtube.Service, query string, start string, end string, max int, pageToken string) *youtube.SearchListCall {
	return service.Search.List([]string{"id, snippet"}).
		Q(query).
		PublishedAfter(start). //"2022-04-07T00:00:00Z").
		PublishedBefore(end).  //"2022-04-08T00:00:00Z").
		MaxResults(int64(max))
}

func callYoutubeAPI(key string, query string, start string, end string, max int) map[string]string {
	service := makeService(key)

	nextPageToken := ""
	responses := [](*youtube.SearchListResponse){}
	// １回のリクエストで取得できるのは25件までなので25件ずつ取得
	for i := 0; i*25 < max; i++ {
		apiCall := makeApiCall(service, query, start, end, i*25, nextPageToken) // TODO 取得件数修正
		response, err := apiCall.Do()
		if err != nil {
			log.Fatal(err)
		}
		responses = append(responses, response)
		nextPageToken = response.NextPageToken
	}

	searchResult := make(map[string]string)
	for _, response := range responses {
		for _, item := range response.Items {
			if item.Id.Kind != "youtube#video" {
				continue
			}
			searchResult[item.Id.VideoId] = item.Snippet.Title
		}
	}
	test := formatIdTitle(searchResult)
	fmt.Println(test)
	return searchResult
}

func formatIdTitle(idTitle map[string]string) []string {
	results := make([]string, 0)
	for id, title := range idTitle {
		results = append(results, fmt.Sprintf("%v\nhttps://www.youtube.com/watch?v=%v\n\n", title, id))
	}
	return results
}
