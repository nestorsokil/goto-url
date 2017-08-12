package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/nestorsokil/gl"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/service"
	"github.com/nestorsokil/goto-url/util"
)

var urlService service.UrlService
var conf util.ApplicationConfig
var logger gl.Logger

func main() {
	conf = util.LoadConfig()

	globalLog := conf.GetGlobalLogFile()
	defer globalLog.Sync()
	defer globalLog.Close()
	requestLog := conf.GetRequestLogFile()
	defer requestLog.Sync()
	defer requestLog.Close()

	ds, err := db.CreateDataSource(&conf)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer ds.Shutdown()

	stop := make(chan struct{})
	defer close(stop)

	logger = gl.Simple().
		WriteTo(globalLog).
		WithLevel(gl.LEVEL_DEBUG).
		WithPrefix("Main").
		Build()
	urlService = service.New(ds, &conf, logger)
	go urlService.ClearRecordsAsync(stop)

	fs := http.FileServer(http.Dir("static"))
	router := mux.NewRouter()
	router.PathPrefix("/home/").Handler(http.StripPrefix("/home/", fs))
	router.Handle("/", http.HandlerFunc(redirectToIndex)).Methods("GET")
	router.Handle("/short", http.HandlerFunc(shorten)).Methods("GET")
	router.Handle("/{key}", http.HandlerFunc(redirect)).Methods("GET")
	withLog := handlers.LoggingHandler(requestLog, router)

	staticFS := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))

	logger.Info("Starting server on %v.", conf.Port)
	http.ListenAndServe(conf.Port, withLog)
}

func shorten(response http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()
	record, err := urlService.GetRecord(
		urlService.RequestBuilder().
			ForUrl(queryParams.Get("url")).
			WithCustomKey(queryParams.Get("customKey")).
			WithCustomExpirationTime(queryParams.Get("customExpire")).
			Build())
	if err != nil {
		logger.Error(err.Error())
		respond(response, http.StatusInternalServerError,
			fmt.Sprintf("Could not shorten URL. Error: %v", err))
		return
	}
	responseUrl := urlService.ConstructUrl(request.URL.Host, record.Key)
	io.WriteString(response, responseUrl)
}

func redirect(response http.ResponseWriter, request *http.Request) {
	key := mux.Vars(request)["key"]
	if key == "" {
		respond(response, http.StatusBadRequest, "No key provided.")
		return
	}
	record, err := urlService.FindByKey(key)
	if err != nil {
		respond(response, http.StatusNotFound, err.Error())
		return
	}
	http.Redirect(response, request, record.URL, http.StatusMovedPermanently)
}

func redirectToIndex(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "home/", http.StatusMovedPermanently)
}

func respond(response http.ResponseWriter, status int, content string) {
	response.WriteHeader(status)
	response.Header().Set("Content-Type", "text/plain")
	io.WriteString(response, content)
}
