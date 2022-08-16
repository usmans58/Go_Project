package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/usmans58/bookings/internal/config"
	handler "github.com/usmans58/bookings/internal/handlers"
	"github.com/usmans58/bookings/internal/render"
)

const portnumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handler.NewRepo(&app)
	handler.NewHandlers(repo)
	render.NewTemplates((&app))

	// http.HandleFunc("/", handler.Repo.HOME)
	// http.HandleFunc("/about", handler.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portnumber))

	// _ = http.ListenAndServe(portnumber, nil)
	serve := &http.Server{
		Addr:    portnumber,
		Handler: routes(&app),
	}
	err = serve.ListenAndServe()
	log.Fatal(err)
}
