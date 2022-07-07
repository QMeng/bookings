package main

import (
	"log"
	"net/http"
	"time"

	"github.com/QMeng/bookings/pkg/config"
	"github.com/QMeng/bookings/pkg/handlers"
	"github.com/QMeng/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false
	render.NewTemplates(&app)

	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
}
