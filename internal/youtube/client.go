package youtube

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YouTubeService wraps the YouTube API client.
type YouTubeService struct {
	service *youtube.Service
}

// NewYouTubeService creates a new YouTubeService instance.
func NewYouTubeService(apiKey string) *YouTubeService {
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	return &YouTubeService{
		service: service,
	}
}

func (s *YouTubeService) ListVideos(query string, maxResults int64) (*youtube.SearchListResponse, error) {
	daysOld := 7

	call := s.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults).
		Type("video").
		PublishedAfter(time.Now().AddDate(0, 0, -daysOld).Format(time.RFC3339))

	return call.Do()
}
