package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/nestorsokil/gl"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/service"
	"github.com/nestorsokil/goto-url/util"
	"github.com/nestorsokil/goto-url/rest"
)

func main() {
	conf := util.LoadConfig()
	globalLog := conf.GetGlobalLogFile()
	defer globalLog.Sync()
	defer globalLog.Close()
	requestLog := conf.GetRequestLogFile()
	defer requestLog.Sync()
	defer requestLog.Close()

	logger := gl.Simple().
		WriteTo(globalLog).
		WithLevel(gl.LEVEL_DEBUG).
		WithPrefix("Main").
		Build()

	ds, err := db.CreateDataSource(&conf)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer ds.Shutdown()

	stop := make(chan struct{})
	defer close(stop)

	urlService := service.New(ds, &conf, logger)
	go urlService.ClearRecordsAsync(stop)

	fs := http.FileServer(http.Dir("static"))
	router := mux.NewRouter()
	router.PathPrefix("/home/").Handler(http.StripPrefix("/home/", fs))
	router.Handle("/", rest.RedirectToIndex()).Methods("GET")
	router.Handle("/{key}", rest.Redirect(urlService)).Methods("GET")
	router.Handle("/url/short", rest.Shorten(urlService)).Methods("GET")
	withLog := handlers.LoggingHandler(requestLog, router)

	logger.Info("Starting server on %v.", conf.Port)
	http.ListenAndServe(conf.Port, withLog)
}
