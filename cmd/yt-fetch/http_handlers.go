package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type httpHandlers struct {
	mongoDB interface {
		get(limit, offset int64) ([]Video, error)
		getQueryMatch(query string, limit, offset int64) ([]Video, error)
	}
}
type PaginatedResponse struct {
	Data  []Video `json:"data"`
	Page  int64   `json:"page"`
	Limit int64   `json:"limit"`
}

func newHTTPHandler(dic *diContainer) (*httpHandlers, error) {
	mongoDB, err := dic.mongodir()
	if err != nil {
		return nil, err
	}
	return &httpHandlers{
		mongoDB: mongoDB,
	}, nil
}

func newHTTPHandlerDIProvider(dic *diContainer) func() (*httpHandlers, error) {
	var s *httpHandlers
	var mu sync.Mutex
	return func() (*httpHandlers, error) {
		mu.Lock()
		defer mu.Unlock()
		var err error
		if s == nil {
			s, err = newHTTPHandler(dic)
		}
		return s, err
	}
}

func (h *httpHandlers) handleGetVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		limit  int64 = 10
		offset int64
	)
	queryParams := r.URL.Query()
	if queryParams.Get("limit") != "" {
		fmt.Println("herefdsvaf")
		limit, _ = strconv.ParseInt(queryParams.Get("limit"), 10, 64)
	}
	if queryParams.Get("offset") != "" {
		fmt.Println("herdeedfrrfrefdsvaf")
		offset, _ = strconv.ParseInt(queryParams.Get("offset"), 10, 64)
	}

	videos, err := h.mongoDB.get(limit, offset)
	if err != nil {
		http.Error(w, "Error getting: "+err.Error(), http.StatusInternalServerError)
	}

	response := PaginatedResponse{
		Data:  videos,
		Page:  offset,
		Limit: limit,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *httpHandlers) handleSearchVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		limit  int64 = 10
		offset int64
	)
	queryParams := r.URL.Query()

	query := queryParams.Get("q")
	if query == "" {
		http.Error(w, "Missing search query parameter", http.StatusBadRequest)
		return
	}
	
	if queryParams.Get("limit") != "" {
		limit, _ = strconv.ParseInt(queryParams.Get("limit"), 10, 64)
	}

	if queryParams.Get("offset") != "" {
		offset, _ = strconv.ParseInt(queryParams.Get("offset"), 10, 64)
	}

	videos, err := h.mongoDB.getQueryMatch(query, limit, offset)
	if err != nil {
		http.Error(w, "failed to get matched query doc", http.StatusInternalServerError)
		return
	}
	response := PaginatedResponse{
		Data:  videos,
		Page:  offset,
		Limit: limit,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response"+err.Error(), http.StatusInternalServerError)
		return
	}
}
