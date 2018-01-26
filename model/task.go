package model

import (
	"sync"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)

type Task struct {
	Id      int     `json:"id"`
	Pages   []*page `json:"pages"`
	IsEnded bool    `json:"-"`
}

type Tasks struct {
	sync    *Throttler
	channel chan *Task
	counter int
	results map[int]*Task
}

func NewTaskStack() *Tasks {
	return &Tasks{
		sync:    NewThrottler(),
		channel: make(chan *Task),
		results: make(map[int]*Task),
	}
}
func (t *Tasks) Add(urls []string) int {
	t.counter++
	t.results[t.counter] = &Task{
		Id:    t.counter,
		Pages: pagesFromUrls(urls),
	}
	go func(task *Task) {
		t.channel <- task
	}(t.results[t.counter])
	return t.counter
}

func (t *Tasks) GetTask(r *http.Request) *Task {
	vars := mux.Vars(r)
	// I don't catch error because route must contain a valid id
	id, _ := strconv.Atoi(vars["id"])
	task, ok := t.results[id]
	if !ok {
		return nil
	}
	return task

}

func (t *Tasks) Crawl() {
	for task := range t.channel {
		go func(task *Task) {
			wg := &sync.WaitGroup{}
			for _, page := range task.Pages {
				wg.Add(1)
				go page.parse(wg, t.sync)
			}
			wg.Wait()
			task.IsEnded = true
		}(task)
	}
}
