package main

import (
	"Balancer/config"
	"Balancer/middleware"
	"Balancer/proxy"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	configure, err := config.ReadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("read configure error: %s", err)
	}
	err = configure.Validation()
	if err != nil {
		log.Fatalf("verify configure error: %s", err)
	}
	router := mux.NewRouter()
	for _, l := range configure.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, l.BalanceMode)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		if configure.HealthCheck {
			httpProxy.HealthCheck(configure.HealthCheckInterval)
		}
		router.Handle(l.Pattern, httpProxy)
	}
	if configure.MaxAllowed > 0 {
		router.Use(middleware.MAxAllowedRequests(configure.MaxAllowed))
	}
	configure.Print()
	svr := http.Server{
		Addr:    ":" + strconv.Itoa(configure.Port),
		Handler: router,
	}

	// listen and serve
	if configure.Schema == "http" {
		err := svr.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	} else if configure.Schema == "https" {
		err := svr.ListenAndServeTLS(configure.SSLCertificate, configure.SSLCertificateKey)
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	}

}
