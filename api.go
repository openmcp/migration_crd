package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "github.com/gorilla/websocket"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	nanumv1alpha1 "nanum.co.kr/openmcp/migration/pkg/apis/nanum/v1alpha1"
	mig "nanum.co.kr/openmcp/migration/pkg/controller/openmcpmigration"
)

// /log
type GetLog struct {
	GetNotSupported
	PostNotSupported
}

func (GetLog) Uri() string {
	return "/log"
}
func (GetLog) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	enableCors(rw)
	print("hihi")

	return Response{200, "", "hihi"}
}

func (GetLog) Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	enableCors(rw)

	var log string
	log = r.FormValue("data")

	// return Response{200, "", migrationsource}

	if log != "" {
		fmt.Print(log)
		return Response{200, "", log}
	} else {
		return Response{400, "", nil}
	}
}

// /migration
type MigrationResource struct {
	GetNotSupported
	PostNotSupported
}

func (MigrationResource) Uri() string {
	return "/migration"
}
func (MigrationResource) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	enableCors(rw)
	return Response{200, "", ""}
}

func (MigrationResource) Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	enableCors(rw)
	var migrationsource nanumv1alpha1.MigrationServiceSource
	json.NewDecoder(r.Body).Decode(&migrationsource)
	// return Response{200, "", migrationsource}

	if migrationsource.VolumePath != "" {
		fmt.Print(migrationsource.MigrationSources[0].TargetCluster)
		mig.MigratioResource(migrationsource.MigrationSources[0])
		return Response{200, "", migrationsource}
	} else {
		return Response{400, "", nil}
	}
}
func enableCors(rw http.ResponseWriter) {

	(rw).Header().Set("Allow", "POST, GET, OPTIONS")
	(rw).Header().Set("Access-Control-Allow-Origin", "*")
	(rw).Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept")
	(rw).Header().Set("Content-Type", "application/json")

}

var addr = flag.String("addr", ":8082", "http service address")

type wsHandler struct {
	Message string `json:"message"`
	Host    string `json:"host"`
	Step    int    `json:"step"`
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// upgrader is needed to upgrade the HTTP Connection to a websocket Connection
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	//Upgrading HTTP Connection to websocket connection
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}
	//handle your websockets with wsConn
	logfile, err := ioutil.ReadFile("logfile.txt")
	wsh = wsHandler{
		Message: string(logfile),
		Host:    "10.0.0.222",
		Step:    1,
	}
	wsConn.WriteJSON(wsh)
}
