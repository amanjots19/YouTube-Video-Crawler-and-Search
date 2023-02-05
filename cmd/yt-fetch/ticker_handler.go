package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/api/youtube/v3"
)

const (
	query        = "tea"
	maxResult    = 50
	developerKey = "AIzaSyBQnJvQC8EXBMk9HsaLpcMo0IulGaBkT_E"
)

var (
	searchPart []string = []string{"id", "snippet"}
	lastdate            = "2023-01-05T12:50:57+05:30"
)

type Video struct {
	ID          string    `bson:"id" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	DateTime    time.Time `bson:"date_time" json:"publishedAt"`
	Thumbnails  thumbnail `bson:"thumbnail" json:"thumbnails"`
}

type thumbnail struct {
	Default struct {
		Url string `json:"url" bson:"url"`
	} `json:"default" bson:"default"`
}

type tickerHandler struct {
	mongoDB interface {
		getLatestDateTime() (*Video, error)
		insert(video Video) error
	}
	ytService *youtube.Service
}

func newTickerHandler(dic *diContainer) (*tickerHandler, error) {
	service, err := dic.ytService()
	if err != nil {
		return nil, err
	}
	mongo, err := dic.mongodir()
	if err != nil {
		return nil, err
	}
	return &tickerHandler{
		ytService: service,
		mongoDB:   mongo,
	}, nil
}

func newTickerHandlerDIProvider(dic *diContainer) func() (*tickerHandler, error) {
	var s *tickerHandler
	var mu sync.Mutex
	return func() (*tickerHandler, error) {
		mu.Lock()
		defer mu.Unlock()
		var err error
		if s == nil {
			s, err = newTickerHandler(dic)
		}
		return s, err
	}
}

func (t *tickerHandler) handle(ctx context.Context) error {
	latestVideoDoc, err := t.mongoDB.getLatestDateTime()
	if err != nil {
		return err
	}
	if !latestVideoDoc.DateTime.IsZero() {
		lastdate = latestVideoDoc.DateTime.Format(time.RFC3339)
		fmt.Println(lastdate, "now inserted in to db")
	}
	call := t.ytService.Search.List(searchPart).
		Q(query).Type("video").
		MaxResults(maxResult).Fields("items(id(videoId),snippet(title,publishedAt,description,thumbnails(default(url))))").PublishedAfter(lastdate)

	response, err := call.Do()
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}
	for _, item := range response.Items {
		datetime, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			return err
		}
		video := Video{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			DateTime:    datetime,
			Thumbnails: thumbnail{
				Default: struct {
					Url string `json:"url" bson:"url"`
				}{
					Url: item.Snippet.Thumbnails.Default.Url,
				},
			},
		}
		err = t.handleDB(ctx, video)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *tickerHandler) handleDB(ctx context.Context, v Video) error {
	err := t.mongoDB.insert(v)
	if err != nil {
		return err
	}
	return nil
}
