package main

import "google.golang.org/api/youtube/v3"

type diContainer struct {
	addr          string
	tickerHandler func() (*tickerHandler, error)
	ytService     func() (*youtube.Service, error)
	mongodir      func() (*mongoDB, error)
	httpHandlers  func() (*httpHandlers, error)
}

func NewDIContainer() (*diContainer, error) {
	dic := &diContainer{
		addr: ":8181",
	}

	dic.tickerHandler = newTickerHandlerDIProvider(dic)
	dic.ytService = newYTClientDIProvider(dic)
	dic.mongodir = newMongoDIProvider(dic)
	dic.httpHandlers = newHTTPHandlerDIProvider(dic)

	return dic, nil
}
