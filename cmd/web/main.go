package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"snippetbox.pethron.me/cmd/config"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	app := &config.Application{
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Lshortfile),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC),
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  routes(app)(),
	}

	app.ErrorLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
