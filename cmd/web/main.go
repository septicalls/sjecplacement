package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"sjecplacement.in/internal/models"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	drives         *models.DriveModel
	roles          *models.RoleModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config

	dsn := "postgres://default:%s@ep-super-pond-22168364-pooler.ap-southeast-1.postgres.vercel-storage.com:5432/verceldb"
	pass := flag.String("pass", "[REDACTED]", "Password for the PostgreSQL Database")

	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	if *pass == "[REDACTED]" {
		flag.Usage()
		errorLog.Fatal("Please provide a password for the DB")
	}

	cfg.dsn = fmt.Sprintf(dsn, *pass)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	gob.Register(models.Drive{})

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		drives:         &models.DriveModel{DB: db},
		roles:          &models.RoleModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
