package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/gorilla/handlers"
	"fmt"
	"io"

	"github.com/nestorsokil/goto-url/config"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/service"
)

func main() {
	globalLog := config.GetGlobalLogFile()
	defer globalLog.Close()

	log.SetOutput(globalLog)
	requestLog := config.GetRequestLogFile()
	defer requestLog.Close()

	session := db.SetupConnection()
	defer session.Close()

	router := mux.NewRouter()
	router.Handle("/short", http.HandlerFunc(shorten)).Methods("GET")
	router.Handle("/{key}", http.HandlerFunc(redirect)).Methods("GET")

	withLog := handlers.LoggingHandler(requestLog, router)

	go service.ClearRecordsAsync()

	log.Printf("[INFO] Starting server on %v.\n", config.Settings.Port)
	fmt.Printf("[INFO] Starting server on %v.\n", config.Settings.Port)
	http.ListenAndServe(config.Settings.Port, withLog)
}

func shorten(response http.ResponseWriter, request *http.Request) {
	url := request.URL.Query().Get("url")
	if url == "" {
		respond(response, http.StatusBadRequest, "No url provided.")
		return
	}

	var base string
	if config.Settings.DevMode == true {
		base = config.Settings.ApplicationUrl
	} else {
		base = request.URL.Host
	}

	record, err := service.GetRecord(url)
	if err != nil {
		log.Println(err)
		respond(response, http.StatusInternalServerError, "Could not shorten URL.")
		return
	}
	key := record.Key
	final := base + "/" + key
	io.WriteString(response, final)
}

func redirect(response http.ResponseWriter, request *http.Request) {
	key := mux.Vars(request)["key"]
	if key == "" {
		respond(response, http.StatusBadRequest, "No key provided.")
		return
	}

	record := db.Find(key)
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