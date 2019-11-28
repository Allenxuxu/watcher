package watcher

import "errors"

var ErrWatcherStopped = errors.New("watcher stopped")

type Watcher struct {
	exit    chan interface{}
	updates chan interface{}
}

func (w *Watcher) Next() (interface{}, error) {
	for {
		select {
		case <-w.exit:
			return nil,ErrWatcherStopped
		case v := <-w.updates:
			return v, nil
		}
	}
}

func (w *Watcher) Stop() error {
	select {
	case <-w.exit:
	default:
		close(w.exit)
	}
	return nil
}
