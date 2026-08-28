package main

import (
	"flag"
	"gestureparticles/internal/api"
	"gestureparticles/internal/flow021"
	"gestureparticles/internal/importer"
	"gestureparticles/internal/registry"
	"gestureparticles/internal/review"
	"gestureparticles/internal/store"
	"log"
	"net/http"
	"os"
)

func main() {
	path := flag.String("db", "gesture.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	db, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	service := registry.New(db, review.New(db), importer.New(db), flow021.New())
	h := api.New(service)
	if os.Getenv("GESTURE_DEMO") == "1" {
		log.Println("demo mode enabled")
	}
	log.Printf("gesture particle service listening on %s", *addr)
	if err := http.ListenAndServe(*addr, h); err != nil {
		log.Fatal(err)
	}
}
