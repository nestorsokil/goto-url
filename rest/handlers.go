package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nestorsokil/goto-url/service"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// Shorten returns an http handler for URL shortening
func Shorten(service service.UrlService) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		queryParams := request.URL.Query()
		log.Debugf("Request to shorten URL '%v' from IP '%v'", queryParams.Get("url"), request.RemoteAddr)
		record, err := service.CreateRecord(queryParams.Get("url"), queryParams.Get("customKey"))
		if err != nil {
			respond(response, http.StatusInternalServerError, fmt.Sprintf("Could not shorten. Error: %v", err))
			return
		}
		responseURL := service.ConstructURL(request.URL.Host, record.Key)
		io.WriteString(response, responseURL)
	}
}

// Redirect returns an http handler that redirects to full URL
func Redirect(service service.UrlService) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		key := mux.Vars(request)["key"]
		log.Debugf("Request for key '%s' received from IP '%v'", key, request.RemoteAddr)
		if key == "" {
			log.Warnf("No key in the request")
			respond(response, http.StatusBadRequest, "No key provided.")
			return
		}
		record, err := service.FindByKey(key)
		if err != nil {
			log.Warnf("Key '%v' was not found", key)
			respond(response, http.StatusNotFound, err.Error())
			return
		}
		log.Debugf("Request for key '%s', redirecting to '%s'", record.Key, record.URL)
		http.Redirect(response, request, record.URL, http.StatusMovedPermanently)
	}
}

func respond(response http.ResponseWriter, status int, content string) {
	response.WriteHeader(status)
	response.Header().Set("Content-Type", "text/plain")
	io.WriteString(response, content)
}
