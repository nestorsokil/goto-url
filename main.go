package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/gorilla/handlers"
	"fmt"
	"io"
)

func main() {
	setupGlobalLog()
	requestLog := getRequestLogFile()
	setupMongoSession()
	defer session.Close()

	router := mux.NewRouter()
	router.Handle("/short", shorten).Methods("GET")
	router.Handle("/{key}", redirect).Methods("GET")

	withLog := handlers.LoggingHandler(requestLog, router)

	log.Printf("[INFO] Starting server on %v.\n", conf.Port)
	fmt.Printf("[INFO] Starting server on %v.\n", conf.Port)
	http.ListenAndServe(conf.Port, withLog)
}

var shorten = http.HandlerFunc(
	func(response http.ResponseWriter, request *http.Request) {
		url := request.URL.Query().Get("url")
		base := request.URL.Host
		fmt.Println(url)
		var key string
		record := findShort(url)
		if record != nil {
			key = record.Key
		} else {
			record, err := createNewRecord(url)
			if err != nil {
				log.Println("[ERROR] Could not save record.", err)
				response.WriteHeader(http.StatusInternalServerError)
				response.Header().Set("Content-Type", "text/plain")
				io.WriteString(response, "Could not shorten URL.")
			}
			key = record.Key
		}
		final := base + "/" + key
		io.WriteString(response, final)
	})



var redirect = http.HandlerFunc(
	func(response http.ResponseWriter, request *http.Request) {
		key := mux.Vars(request)["key"]
		fmt.Println(key)
		if key == "" {
			response.WriteHeader(http.StatusBadRequest)
			io.WriteString(response, "No key provided.")
			return
		}

		record := find(key)
		if record == nil {
			response.WriteHeader(http.StatusNotFound)
			io.WriteString(response, "URL not found.")
			return
		}

		http.Redirect(response, request, record.URL, http.StatusMovedPermanently)
	})