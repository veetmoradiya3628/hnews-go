package main

import (
	"net/http"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
		// ReadTimeout: 2 * time.Second, // Set a read timeout to prevent slow clients from hanging the server
	}
	return srv.ListenAndServe()
}
