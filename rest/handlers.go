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
		url := queryParams.Get("url")
		log.Debugf("Request to shorten URL '%v' from IP '%v'", url, request.RemoteAddr)
		record, err := service.CreateRecord(request.Context(), url, queryParams.Get("customKey"))
		if err != nil {
			respond(request, response, http.StatusInternalServerError, fmt.Sprintf("Could not shorten. Error: %v", err))
			return
		}
		responseURL := constructURL(request.URL.Host, record.Key)
		respond(request, response, http.StatusCreated, responseURL)
	}
}

// Redirect returns an http handler that redirects to full URL
func Redirect(service service.UrlService) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		key := mux.Vars(request)["key"]
		log.Debugf("Request for key '%s' received from IP '%v'", key, request.RemoteAddr)
		if key == "" {
			log.Warnf("No key in the request")
			respond(request, response, http.StatusBadRequest, "No key provided.")
			return
		}
		record, err := service.FindByKey(request.Context(), key)
		if err != nil {
			log.Warnf("Key '%v' was not found", key)
			respond(request, response, http.StatusNotFound, err.Error())
			return
		}
		log.Debugf("Request for key '%s', redirecting to '%s'", record.Key, record.URL)
		http.Redirect(response, request, record.URL, http.StatusMovedPermanently)
	}
}

func constructURL(host, key string) string {
	return host + "/" + key
}

func respond(req *http.Request, res http.ResponseWriter, status int, content string) {
	res.WriteHeader(status)
	res.Header().Set("Content-Type", "text/plain")
	if _, err := io.WriteString(res, content); err != nil {
		log.Errorf("Could not write response to IP %v. Error: %v", req.RemoteAddr, err.Error())
	}
}
