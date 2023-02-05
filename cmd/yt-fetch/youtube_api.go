package main

import (
	"context"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func newYTClient(dic *diContainer) (*youtube.Service, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(developerKey))
	if err != nil {
		return nil, err
	}
	return service, nil
}

func newYTClientDIProvider(dic *diContainer) func() (*youtube.Service, error) {
	var s *youtube.Service
	var mu sync.Mutex
	return func() (*youtube.Service, error) {
		mu.Lock()
		defer mu.Unlock()
		var err error
		if s == nil {
			s, err = newYTClient(dic)
		}
		return s, err
	}
}
