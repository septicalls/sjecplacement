package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
	}

	srv := http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", cfg.addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
