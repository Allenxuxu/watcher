package watcher

import (
	"fmt"
	"sync"
	"time"
)

func ExampleNewSource() {
	source := NewSource()
	w1, err := source.Watch()
	if err != nil {
		panic(err)
	}
	w2, err := source.Watch()
	if err != nil {
		panic(err)
	}

	if source.watchers.Len() != 2 {
		panic("watchers's len should be 2")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		v, err := w1.Next()
		if err != nil {
			panic(err)
		}

		fmt.Println(v)
		wg.Done()
	}()
	go func() {
		if err := w2.Stop(); err != nil && err != ErrWatcherStopped {
			panic(err)
		}
		_, err := w2.Next()
		if err == nil {
			panic(err)
		}

		wg.Done()
	}()

	time.Sleep(time.Second)
	source.Update("hello")

	wg.Wait()

	if source.watchers.Len() != 1 {
		panic("watchers's len should be 1")
	}

	if err := w1.Stop(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	if source.watchers.Len() != 0 {
		panic("watchers's len should be 0")
	}

	// output:
	// hello

}
