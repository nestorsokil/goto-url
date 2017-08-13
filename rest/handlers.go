package rest

import (
	"net/http"
	"fmt"
	"io"
	"github.com/gorilla/mux"
	"github.com/nestorsokil/goto-url/service"
)

func Shorten(service service.UrlService) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		queryParams := request.URL.Query()
		record, err := service.GetRecord(
			service.RequestBuilder().
				ForUrl(queryParams.Get("url")).
				WithCustomKey(queryParams.Get("customKey")).
				WithCustomExpirationTime(queryParams.Get("customExpire")).
				Build())
		if err != nil {
			respond(response, http.StatusInternalServerError,
				fmt.Sprintf("Could not shorten URL. Error: %v", err))
			return
		}
		responseUrl := service.ConstructUrl(request.URL.Host, record.Key)
		io.WriteString(response, responseUrl)
	}
}

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
		http.Redirect(response, request, record.URL, http.StatusMovedPermanently)
	}
}

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