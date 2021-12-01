package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/khatibomar/goction/item"
)

type application struct {
	// TODO(khatibomar) : make it a slice of items
	host     string
	item     *item.Item
	errorLog *log.Logger
	infoLog  *log.Logger
}

type srv struct {
	server *http.Server
	raft   *raft.Raft
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	potatoItem := item.New("potato", 50, "$")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		host:     "localhost" + *addr,
		item:     potatoItem,
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := srv{
		server: &http.Server{
			Addr:         *addr,
			ErrorLog:     errorLog,
			Handler:      app.routes(),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.server.ListenAndServe()
	errorLog.Fatal(err)
}
