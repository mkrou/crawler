package server

import (
	"net/http"
	"bufio"
	"github.com/mkrou/crawler/model"
)

func SetTask(tasks *model.Tasks) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		scanner := bufio.NewScanner(r.Body)
		var urls []string
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			JsonMsgResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(urls) == 0 {
			JsonMsgResponse(w, "List of urls is empty.", http.StatusBadRequest)
			return
		}
		JsonResponse(w, CreateResponse{tasks.Add(urls)}, http.StatusCreated)
	}
}
func GetTask(tasks *model.Tasks) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var task *model.Task
		if task = tasks.GetTask(r); task == nil {
			JsonMsgResponse(w, "Task isn't found.", http.StatusNotFound)
			return
		}
		if !task.IsEnded {
			JsonMsgResponse(w, "Task isn't finished.", http.StatusNoContent)
			return
		}
		JsonResponse(w, task, http.StatusOK)
	}
}
