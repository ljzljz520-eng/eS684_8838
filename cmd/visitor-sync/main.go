package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"parkvisitor/internal/clock"
	"parkvisitor/internal/service"
	"parkvisitor/internal/storage"
	"parkvisitor/internal/transport"
)

func main() {
	path := flag.String("db", "visitor-sync.db", "bbolt database path")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	store, err := storage.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	app, err := service.NewApp(store, clock.Fixed())
	if err != nil {
		log.Fatal(err)
	}
	handler := transport.NewHandler(app)
	fmt.Printf("park visitor sync listening on %s using %s\n", *addr, *path)
	if err := http.ListenAndServe(*addr, handler); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
