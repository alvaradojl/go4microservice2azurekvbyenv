package main

import (
	log "log"
	"net/http"

	"github.com/opentracing-contrib/go-stdLib/nethttp"

	opentracing "github.com/opentracing/opentracing-go"
)

const (
	listenAddress = ":8080"
)

var httpClient = &http.Client{Transport: &nethttp.Transport{}} //TODO: refactor variables for better visibility

func main() {

	log.Println("starting keyvault microservice app...")

	tracer, closer := initJaeger("go4microservice2azurekvbyenv")

	opentracing.SetGlobalTracer(tracer)

	defer closer.Close()

	http.Handle("/", http.HandlerFunc(getSecret))

	mux := nethttp.Middleware(tracer, http.DefaultServeMux)

	err := http.ListenAndServe(listenAddress, mux)

	if err != nil {
		log.Fatal(err)
	}

}
