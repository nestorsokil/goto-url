package rest

import (
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/nestorsokil/goto-url/service"
	log "github.com/sirupsen/logrus"
)

// Shorten returns an http handler for URL shorening
func Shorten(service service.UrlService) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		queryParams := request.URL.Query()
		record, err := service.GetRecord(
			service.RequestBuilder().
				ForURL(queryParams.Get("url")).
				WithCustomKey(queryParams.Get("customKey")).
				WithCustomExpirationTime(queryParams.Get("customExpire")).
				Build())
		if err != nil {
			respond(response, http.StatusInternalServerError,
				fmt.Sprintf("Could not shorten URL. Error: %v", err))
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
		if key == "" {
			respond(response, http.StatusBadRequest, "No key provided.")
			return
		}
		record, err := service.FindByKey(key)
		if err != nil {
			respond(response, http.StatusNotFound, err.Error())
			return
		}
		log.Debugf("Request for key '%s', redirecting to '%s'", record.Key, record.URL)
		http.Redirect(response, request, record.URL, http.StatusMovedPermanently)
	}
}

// RedirectToIndex returns an http handler taht sends user to landing page
func RedirectToIndex() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "home/", http.StatusMovedPermanently)
	}
}

func respond(response http.ResponseWriter, status int, content string) {
	response.WriteHeader(status)
	response.Header().Set("Content-Type", "text/plain")
	io.WriteString(response, content)
}
