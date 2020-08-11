package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Resource interface {
	Uri() string
	Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
}

type (
	GetNotSupported  struct{}
	PostNotSupported struct{}
)

func (GetNotSupported) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "get not supported", nil}
}
func (PostNotSupported) Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "post not supported", nil}
}
func abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}
func HttpResponse(rw http.ResponseWriter, req *http.Request, res Response) {
	content, err := json.Marshal(res)

	if err != nil {
		abort(rw, 500)
	}
	rw.WriteHeader(res.Code)
	rw.Write(content)
}

func AddResource(router *httprouter.Router, resource Resource) {
	fmt.Println("\"" + resource.Uri() + "\" api is registerd")

	router.GET(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Get(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.POST(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Post(rw, r, ps)
		HttpResponse(rw, r, res)
	})
}

//open websocket
func server() {
	router := http.NewServeMux()
	router.Handle("/ws", wsHandler{})
	log.Fatal(http.ListenAndServe("10.0.0.222:8083", router))
}

func main() {

	router := httprouter.New()
	AddResource(router, new(MigrationResource))
	AddResource(router, new(GetLog))
	go server()
	log.Fatal(http.ListenAndServe("10.0.0.222:8082", router))

}
