package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/rest"
	"github.com/nestorsokil/goto-url/service"
	"github.com/nestorsokil/goto-url/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	conf := util.LoadConfig()
	requestLog := conf.GetRequestLogFile()
	defer requestLog.Sync()
	defer requestLog.Close()

	ds, err := db.CreateDataSource(&conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ds.Shutdown()

	stop := make(chan struct{})
	defer close(stop)

	urlService := service.New(ds, &conf)
	go urlService.ClearRecordsAsync(stop)

	fs := http.FileServer(http.Dir("static"))
	router := mux.NewRouter()
	router.PathPrefix("/home/").Handler(http.StripPrefix("/home/", fs))
	router.Handle("/", rest.RedirectToIndex()).Methods("GET")
	router.Handle("/{key}", rest.Redirect(urlService)).Methods("GET")
	router.Handle("/url/short", rest.Shorten(urlService)).Methods("GET")
	withLog := handlers.LoggingHandler(requestLog, router)

	log.Infof("Starting server on %v.", conf.Port)
	http.ListenAndServe(conf.Port, withLog)
}
