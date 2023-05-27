package main

import (
	"net/http"
	"time"
)

func main() {
	router := initRouter()
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
