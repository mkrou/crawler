package server

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mkrou/crawler/model"
)



type CreateResponse struct {
	Id int `json:"id"'`
}
type MessageResponse struct {
	Message string `json:"message"`
}

func JsonMsgResponse(w http.ResponseWriter, msg string, statusCode int) error {
	return JsonResponse(w, MessageResponse{msg}, statusCode)
}
func JsonResponse(w http.ResponseWriter, object interface{}, statusCode int) error {
	resp, err := json.Marshal(object)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s\n", resp)
	return nil
}
func Serve(tasks *model.Tasks)error{
	r := mux.NewRouter()
	r.HandleFunc("/tasks", SetTask(tasks)).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{id:[0-9]+}", GetTask(tasks)).Methods(http.MethodGet)
	return http.ListenAndServe(":8080", r)
}