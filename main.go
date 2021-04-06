package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "Google", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
)

//func getChannel(client *http.Client)

func main() {

	apiKey, err := ioutil.ReadFile("./secrets/youtube_api_key")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(string(apiKey)))
	if err != nil {
		log.Fatalf("error creating youtube client: %v", err)
	}
	call := service.Search.List([]string{"id,snippet"}).Q("cat").MaxResults(25)
	response, err := call.Do()
	//fmt.Println(response)
	if err != nil {
		log.Fatal("bad call")
	}

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	printIDs("Videos", videos)
	printIDs("Channels", channels)
	printIDs("Playlists", playlists)

	fmt.Println("yes")
	http.Handle("/", http.FileServer(http.Dir("./views")))
	http.ListenAndServe(":3000", nil)
}

func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}
