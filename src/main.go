package main

import (
	"access-tester/src/scraper"
	"access-tester/src/server"
	"github.com/go-chi/chi/v5"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	log.Printf("[main] initializing...")
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "localhost"
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	swagger, err := server.GetSwagger()
	if err != nil {
		panic(err)
	}

	swagger.Servers = nil

	handler := server.NewStrictHandler(server.ServerImpl{Scraper: scraper.NewScraper()}, nil)
	router := chi.NewRouter()
	router.Use(nethttpmiddleware.OapiRequestValidator(swagger))
	server.HandlerFromMux(handler, router)
	srv := http.Server{
		Handler: router,
		Addr:    net.JoinHostPort(host, port),
	}
	log.Println("[main] initialization complete")
	log.Printf("[main] serving on %s:%s", host, port)
	log.Fatal(srv.ListenAndServe())
}
