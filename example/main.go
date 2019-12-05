package main

import (
	"log"
	"sync"
	"time"

	"github.com/Allenxuxu/watcher"
)

func main() {
	source := watcher.NewSource()
	w1, err := source.Watch()
	if err != nil {
		panic(err)
	}
	w2, err := source.Watch()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		v, err := w1.Next()
		if err != nil {
			panic(err)
		}

		log.Println("1 :", v)
		wg.Done()
	}()
	go func() {
		if err := w2.Stop(); err != nil && err != watcher.ErrWatcherStopped {
			panic(err)
		}
		v, err := w2.Next()
		if err != nil && err != watcher.ErrWatcherStopped {
			panic(err)
		}

		log.Println("2 :", v)
		wg.Done()
	}()

	time.Sleep(time.Second)
	source.Update("hello")

	wg.Wait()
}
