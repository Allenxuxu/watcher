package watcher

import (
	"container/list"
	"sync"
)

type Source struct {
	sync.RWMutex
	value interface{}

	watchers *list.List
}

func NewSource() *Source {
	return &Source{
		watchers: list.New(),
	}
}

func (s *Source) Read() (interface{}, error) {
	s.RLock()
	v := s.value
	s.RUnlock()
	return v, nil
}

func (s *Source) Watch() (*Watcher, error) {
	w := &Watcher{
		exit:    make(chan interface{}),
		updates: make(chan interface{}, 1),
	}

	s.Lock()
	element := s.watchers.PushBack(w)
	s.Unlock()

	go func() {
		<-w.exit
		s.Lock()
		s.watchers.Remove(element)
		s.Unlock()
	}()

	return w, nil
}

func (s *Source) Update(v interface{}) {
	s.Lock()
	s.value = v
	s.Unlock()

	watchers := s.copyWatchers()
	for _, w := range watchers {
		// 防止阻塞, 如果发送时 chan 缓冲区满了，则丢次本次 update 内容
		select {
		case w.updates <- v:
		default:
		}
	}
}

// BUpdate blocking update
func (s *Source) BUpdate(v interface{}) {
	s.Lock()
	s.value = v
	s.Unlock()

	watchers := s.copyWatchers()
	for _, w := range watchers {
		w.updates <- v
	}
}

func (s *Source) copyWatchers() []*Watcher {
	watchers := make([]*Watcher, 0, s.watchers.Len())

	s.RLock()
	for e := s.watchers.Front(); e != nil; e = e.Next() {
		watchers = append(watchers, e.Value.(*Watcher))
	}
	s.RUnlock()
	return watchers
}
