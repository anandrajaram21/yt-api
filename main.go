package main

import (
	"log"
	"time"

	db "github.com/anandrajaram21/yt-api/database"
	models "github.com/anandrajaram21/yt-api/models"
)

func main() {
	database := db.SetupDatabase()
	service, err := initYoutubeClient()
	listArr := []string{"id", "snippet"}

	if err != nil {
		log.Fatalf("Youtube Client could not be initialized")
	}

	searchTerm := "football"

	daysOld := 7

	call := service.Search.List(listArr).
		Q(searchTerm).
		MaxResults(20).
		Type("video").
		PublishedAfter(time.Now().AddDate(0, 0, -daysOld).Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error calling Youtube API: %v", err)
	}

	for _, item := range response.Items {
		log.Println(item.Snippet.Title)
	}

	for _, item := range response.Items {
		publishedAtTime, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Fatalf("Error parsing time: %v", err)
		}

		video := models.Video{
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishDate:  publishedAtTime,
			ThumbnailURL: item.Snippet.Thumbnails.Default.Url,
		}

		result := database.Create(&video)

		if result.Error != nil {
			log.Fatalf("Error saving video: %v", result.Error)
		}

	}
}
