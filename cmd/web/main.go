package main

import (
	"com.snippetbox.aitu/internal/models"
	"context"
	"flag"
	"github.com/jackc/pgx/v5"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// DATABASE CONNECTION
	DBS := "postgres://postgres:1234567@localhost:5432/snippetbox"
	pool, error := pgx.Connect(context.Background(), DBS)
	if error != nil {
		errorLog.Fatalf("Unable to connection to database: %v\n", error)
	}
	defer pool.Close(context.Background())

	templateCache, errr := newTemplateCache()
	if errr != nil {
		errorLog.Fatal(errr)
	}

	// testing
	//var title string
	//var content string
	//error = pool.QueryRow(context.Background(), "select title, content from snippets where id = 2").Scan(&title, &content)
	//if error != nil {
	//	errorLog.Printf("QueryRow failed: %v\n", error)
	//	os.Exit(1)
	//}
	//
	//fmt.Println(title)

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: pool},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
