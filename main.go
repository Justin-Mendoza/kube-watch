package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
)

// helper function
func healthCheckHandler (w http.ResponseWriter, r *http.Request){
	// Fprintf send sinmple response to client
	fmt.Fprintf(w, "Service is healthy and running on port 8080!")
}

// The main fn, entry point of program exec, the main

func main(){
	// http.HandleFunc registers handler func to URL path : mapping root to healthCheckFn
	// "/" will be root path: localhost:8080/
	// when request are done to this path, we call healthCheckHandler
	http.HandleFunc("/", healthCheckHandler)

	// os.Getenv gets value from env named "PORT"
	// env variables are disjointed from main for security reasons

	port := os.Getenv("PORT")
	// := is for declaraing and assignment

	// we want to check if the port is empty:
	if port == "" {
		// set default to 8080
		port= "8080"
	}

	// log.printf( to print string to log/terminal)
	log.Printf("Starting Server on :%s", port)

	// http.ListenAndServe, starts http server, 2 arguments: 
	// 1. address to listen to
	// 2. handler fn ( nil to use fault set by HandleFunc)
	err := http.ListenAndServe(":"+port, nil)


	// if http.ListenAndServe returns error, block execution 
	if err != nil{
		// log.Fatalf error prints faral error message, and exits program w status 1
		log.Fatalf("Server Failed to Start: %v", err)
	}


}		