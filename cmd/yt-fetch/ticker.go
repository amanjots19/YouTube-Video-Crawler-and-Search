package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const interval = 10 * time.Second

func ticker(ctx context.Context, dic *diContainer, wg *sync.WaitGroup) {
	fmt.Println("herre")
	handler, err := dic.tickerHandler()
	if err != nil {
		log.Println("Error getting TickerHandler:", err)
		wg.Done()
		return
	}
	ticker := time.NewTicker(interval)

	errCh := make(chan error, 1)
	go func() {
		for range ticker.C {
			fmt.Println("Handling Ticker")
			ctx := context.Background()
			select {
			case <-ticker.C:
				err = handler.handle(ctx)
				if err != nil {
					fmt.Println(err)
					errCh <- err
					return
				}
			case <-ctx.Done():
				return
			}
		}
		wg.Done()
	}()

	go func() {
		for err := range errCh {
			log.Println("Error in Ticker:", err)
		}
	}()
}
