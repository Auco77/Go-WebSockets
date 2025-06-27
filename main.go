package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//Create a root ctx and a CancelFunc which can be used to cancel retentionMap goroutine
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)

	defer cancel()

	setupApi(ctx)
	log.Fatal(http.ListenAndServe(":8080", nil))
	println("http server listen")
}

func setupApi(ctx context.Context) {
	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/login", manager.loginHandler)
	http.HandleFunc("/ws", manager.serveWS)

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, len(manager.clients))
	})
}
