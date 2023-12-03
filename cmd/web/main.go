package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"sjecplacement.in/internal/models"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	addr string
	env  string
	db   struct {
		dsn string
	}
	staticDir string
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
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	if err := godotenv.Load(); err != nil {
		errorLog.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:      os.Getenv("APPLICATION_PORT"),
		staticDir: "./ui/static",
	}

	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	cfg.db.dsn = os.Getenv("POSTGRES_URL")

	db, err := openDB(cfg)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	infoLog.Printf("database connection pool established")

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

	app.infoLog.Printf("starting %s server on %s", cfg.env, cfg.addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
