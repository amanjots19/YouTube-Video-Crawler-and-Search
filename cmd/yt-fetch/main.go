package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)
func main() {
	log.Println("----- application started -----")
	dic, err := NewDIContainer()
	if err != nil {
		log.Fatal("get dependencies")
	}
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg.Add(2)
	go ticker(ctx, dic, &wg)

	go func() {
		fmt.Println("heufeoijfknoerf")
		err = runHTTPServer(ctx, dic)
		if err != nil {
			log.Fatal(err)
			cancel()
		}
		wg.Done()
	}()
	wg.Wait()
	<-ctx.Done()
	log.Println("----- application stopped -----")
}
