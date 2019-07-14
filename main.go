package main

import (
	"github.com/gorilla/mux"
	"github.com/nestorsokil/goto-url/conf"
	"net/http"

	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/rest"
	"github.com/nestorsokil/goto-url/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.DebugLevel)

	var c conf.Config = &conf.EnvConfig{}
	ds, err := db.CreateStorage(c)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer ds.Shutdown()
	urlService := service.New(ds, c)
	router := mux.NewRouter()
	router.Handle("/{key}", rest.Redirect(urlService)).Methods("GET")
	router.Handle("/url/short", rest.Shorten(urlService)).Methods("GET")

	port := c.GetString(conf.EnvPort)
	log.Infof("Starting server on %v.", port)
	log.Error(http.ListenAndServe(port, router))
}
