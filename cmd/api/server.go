package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	srv := http.Server{
		Addr:         fmt.Sprintf(":%v", app.cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	err := srv.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}
