package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	listen  = flag.String("listen", ":10000", "listen address")
	datadir = flag.String("datadir", "data/", "data directory")
	maxsize = flag.Int64("maxsize", 10*1024*1024, "maximaum file size")
	baseurl = flag.String("baseurl", "http://localhost:10000", "baseurl")
)

func main() {
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir(*datadir))))
	http.Handle("/upload.cgi", &UploadHandler{
		DataDir:     *datadir,
		MaxFileSize: *maxsize,
		BaseURL:     *baseurl,
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	srv := http.Server{
		Addr: *listen,
	}
	go func() {
		log.Printf("SIGNAL %d received, then shutting down...\n", <-sigCh)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Print(err)
		}
		log.Print("Server shutdown")
	}()

	log.Printf("Listen: %s", *listen)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Server closed with error:", err)
	}
	log.Print("finished")
}
