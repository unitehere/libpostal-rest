package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
)

// Request format for single query requests
type Request struct {
	Queries []string `json:"queries"`
}

func main() {
	host := os.Getenv("LISTEN_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "8080"
	}
	listenSpec := fmt.Sprintf("%s:%s", host, port)

	router := mux.NewRouter()
	router.HandleFunc("/health", HealthHandler).Methods("GET")
	router.HandleFunc("/expand", ExpandHandler).Methods("POST")
	router.HandleFunc("/parser", ParserHandler).Methods("POST")

	s := &http.Server{Addr: listenSpec, Handler: router}
	go func() {
		fmt.Printf("listening on http://%s\n", listenSpec)
		s.ListenAndServe()
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	<-stop
	fmt.Println("\nShutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	s.Shutdown(ctx)
	fmt.Println("Server stopped")
}

// HealthHandler always returns an OK response if the service is healthy
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ExpandHandler calls the libpostal expand functionality
func ExpandHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req Request

	q, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(q, &req)

	result := make([][]string, len(req.Queries))

	for i := 0; i < len(req.Queries); i++ {
		expansions := expand.ExpandAddress(req.Queries[i])
		result[i] = expansions
	}

	serializedResult, _ := json.Marshal(result)
	w.Write(serializedResult)
}

// ParserHandler calls the libpostal parse functionality
func ParserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req Request

	q, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(q, &req)

	result := make([][]parser.ParsedComponent, len(req.Queries))

	for i := 0; i < len(req.Queries); i++ {
		parsed := parser.ParseAddress(req.Queries[i])
		result[i] = parsed
	}

	serializedResult, _ := json.Marshal(result)
	w.Write(serializedResult)
}
