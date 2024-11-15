package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Use otelhttp to record a span for the outgoing call, and propagate
	// context to the destination.
	destination := os.Getenv("DESTINATION_URL")
	resp, err := http.Get(destination)
	if err != nil {
		log.Fatal("could not fetch remote endpoint")
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("could not read response from %v", destination)
	}

	fmt.Fprint(w, strconv.Itoa(resp.StatusCode))
}

func main() {
	// Handle incoming request.
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	var handler http.Handler = r

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), handler))
}
