package youtube

import (
	"context"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var youtubeService *youtube.Service

func InitializeYouTubeAPI(apiKey string) error {
	var err error
	youtubeService, err = youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	return nil
}

func ListVideos(query string, maxResults int64) (*youtube.SearchListResponse, error) {
	daysOld := 7

	call := youtubeService.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults).
		Type("video").
		PublishedAfter(time.Now().AddDate(0, 0, -daysOld).Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}
