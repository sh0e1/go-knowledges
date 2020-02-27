package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	cctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-sigch
		log.Printf("received signal: %s\n", sig.String())
		cancel()
	}()

	run(cctx)
}

func run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	var wg sync.WaitGroup

label:
	for {
		select {
		case <-ctx.Done():
			log.Printf("%v", ctx.Err())
			break label
		case <-ticker.C:
			vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
			for _, v := range vals {
				wg.Add(1)
				go work(ctx, &wg, v)
			}
		}
	}

	wg.Wait()
}

func work(ctx context.Context, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	fmt.Printf("start: %d\n", i)
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Printf("end: %d\n", i)
}
