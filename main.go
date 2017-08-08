package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"

	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/service"
	"github.com/nestorsokil/goto-url/util"
)

var urlService service.UrlService
var conf util.Configuration

func main() {
	conf = util.LoadConfig()

	globalLog := conf.GetGlobalLogFile()
	defer globalLog.Close()
	log.SetOutput(globalLog)

	requestLog := conf.GetRequestLogFile()
	defer requestLog.Close()

	session := db.NewMongoSession(&conf)
	defer session.Close()

	stop := make(chan struct{})
	defer close(stop)

	urlService = service.New(
		db.NewMongoDS(session, conf.Database), &conf)
	go urlService.ClearRecordsAsync(stop)

	router := mux.NewRouter()
	router.Handle("/short", http.HandlerFunc(shorten)).Methods("GET")
	router.Handle("/{key}", http.HandlerFunc(redirect)).Methods("GET")
	withLog := handlers.LoggingHandler(requestLog, router)

	log.Printf("[INFO] Starting server on %v.\n", conf.Port)
	fmt.Printf("[INFO] Starting server on %v.\n", conf.Port)
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
		log.Println(err)
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
	record := urlService.FindByKey(key)
	if record == nil {
		respond(response, http.StatusNotFound, "URL not found.")
		return
	}
	http.Redirect(response, request, record.URL, http.StatusMovedPermanently)
}

func respond(response http.ResponseWriter, status int, content string) {
	response.WriteHeader(status)
	response.Header().Set("Content-Type", "text/plain")
	io.WriteString(response, content)
}
