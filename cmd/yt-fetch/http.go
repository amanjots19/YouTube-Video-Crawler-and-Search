package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func runHTTPServer(ctx context.Context, dic *diContainer) error {
	fmt.Println("hersssssssse")
	h, err := dic.httpHandlers()
	if err != nil{
		return err
	}
	router := mux.NewRouter()
	registerHandlers(router, h)
	fmt.Println("dabhcdsnbkjdnds")
	err = http.ListenAndServe(dic.addr, router)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
	<-ctx.Done()
	return nil
}

// registerHandlers registers the handle functions in the routes.
func registerHandlers(router *mux.Router, h *httpHandlers) {
	router.HandleFunc("/videos", h.handleGetVideos).Methods(http.MethodGet)
	router.HandleFunc("/search", h.handleSearchVideos).Methods(http.MethodGet)
}
