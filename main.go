package main

import (
	"banknote-tracker-auth/handler"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var logger = log.New(os.Stdout, "==>   ", log.LstdFlags)

func main() {
	router := mux.NewRouter()
	handler := handler.New(logger)

	handler.SetupRoutes(router)

	fmt.Println("serving on port 2000...")

	server := &http.Server{
		Addr:         ":2000",
		Handler:      router,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		IdleTimeout:  time.Second * 60,
	}

	server.ListenAndServe()
}
