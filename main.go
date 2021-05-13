package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	//"net/http"

	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

type Liver struct {
	Name string `json:"name,omniempty"`

	Slug string `json:"slug,omniempty"`

	Affiliation string `json:"affiliation,omniempty"`

	EnglishName string `json:"english_name,omniempty"`

	YoutubeURL string `json:"youtube_ch,omniempty"`

	TwitterURL string `json:"twitter,omniempty"`
}

type Video struct {
	Title string

	URL string
}

type Service struct {
	*youtube.Service
}

var (
	query      = flag.String("query", "Google", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
)

//func getChannel(client *http.Client)

func makeLiverData() map[string]Liver {

	liverData := make(map[string]Liver)

	file, err := os.Open("livers.txt")
	if err != nil {
		log.Fatalf("failed to open file")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	for _, line := range lines {
		var liver Liver

		info := strings.SplitN(line, ":", 2)
		err = json.Unmarshal([]byte(info[1]), &liver)
		if err != nil {
			fmt.Println(err)
		}
		liverData[info[0]] = liver
	}

	return liverData
}

func getYoutubeChannelID(liver Liver) string {
	if liver.YoutubeURL == "" {
		return ""
	}
	ch := strings.Split(liver.YoutubeURL, "/")
	return ch[len(ch)-1]
}

func initializeYoutubeService() *Service {
	apiKey, err := os.ReadFile("./secrets/youtube_api_key")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(string(apiKey)))
	if err != nil {
		log.Fatalf("error creating youtube client: %v", err)
	}
	return &Service{service}
}

// TODO figure out how to use less calls to API (maybe use eventType parameter to call once then search)
func (s *Service) getUpcomingStreams(liverData map[string]Liver) map[string][]Video {
	upcomingStreams := make(map[string][]Video)

	for liver, data := range liverData {
		channelId := getYoutubeChannelID(data)
		if channelId != "" {
			call := s.Search.List([]string{"id,snippet"}).ChannelId(channelId).Order("date").MaxResults(50)
			response, err := call.Do()
			if err != nil {
				log.Fatal(err)
			}

			for _, item := range response.Items {
				if item.Id.Kind == "youtube#video" && item.Snippet.LiveBroadcastContent == "upcoming" {
					video := Video{
						Title: item.Snippet.Title,
						URL:   makeYoutubeVideoURL(item.Id.VideoId),
					}
					if streams, ok := upcomingStreams[liver]; ok {
						upcomingStreams[liver] = append(streams, video)
					} else {
						upcomingStreams[liver] = []Video{video}
					}
				}
			}

		}

	}
	return upcomingStreams
}

//func (s *Service) findVideosByChannelId()

func makeYoutubeVideoURL(videoId string) string {
	return "https://www.youtube.com/watch?v=" + videoId
}

func youtubeTest() {

	liverData := makeLiverData()

	service := initializeYoutubeService()

	streams := service.getUpcomingStreams(liverData)
	printStreams(streams)

}

func main() {

	youtubeTest()

	//http.Handle("/", http.FileServer(http.Dir("./views")))
	//http.ListenAndServe(":3000", nil)
}

func printStreams(streams map[string][]Video) {
	fmt.Printf("%v:\n\n", "Upcoming")
	for name, streams := range streams {
		fmt.Printf("%v:\n", name)
		for _, stream := range streams {
			fmt.Printf("Title: %v URL: %v\n", stream.Title, stream.URL)
		}
		fmt.Println()
	}
	fmt.Printf("\n\n")
}
