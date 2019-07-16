package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nestorsokil/goto-url/conf"
	"github.com/nestorsokil/goto-url/db"
	"github.com/nestorsokil/goto-url/rest"
	"github.com/nestorsokil/goto-url/service"
	"github.com/nestorsokil/goto-url/util"
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConf "github.com/uber/jaeger-client-go/config"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.DebugLevel)

	var c conf.Config = &conf.EnvConfig{}
	ds, err := db.CreateStorage(c)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer ds.Shutdown(context.Background())
	urlService := service.New(ds, c)
	router := mux.NewRouter()

	tracer, closer := createTracer()
	defer func() {
		if err := closer.Close(); err != nil {
			log.Errorf("Could not close tracer: %v", err.Error())
		}
	}()
	opentracing.SetGlobalTracer(tracer)

	router.Handle("/metrics", promhttp.Handler())

	router.Handle("/{key}", util.HttpSpan("http.key.redirect", rest.Redirect(urlService))).Methods("GET")
	router.Handle("/url/short", util.HttpSpan("http.key.shorten", rest.Shorten(urlService))).Methods("GET")

	port := c.GetString(conf.EnvPort)
	log.Infof("Starting server on %v.", port)
	log.Error(http.ListenAndServe(port, router))
}

func createTracer() (opentracing.Tracer, io.Closer) {
	cfg, err := jaegerConf.FromEnv()
	if err != nil {
		log.Fatalf("ERROR: cannot init Jaeger: %v\n", err)
	}
	tracer, closer, err := cfg.NewTracer(jaegerConf.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("ERROR: cannot init Jaeger: %v\n", err)
	}
	return tracer, closer
}
