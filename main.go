package main

import (
	"github/somyaranjan99/basic-go-project/cmd/web/middleware"
	"github/somyaranjan99/basic-go-project/pkg/config"
	"github/somyaranjan99/basic-go-project/pkg/handlers"
	"github/somyaranjan99/basic-go-project/pkg/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
)

var app config.AppConfig

func main() {
	app := config.AppConfig{}
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true
	sessionManager := scs.New()
	sessionManager.Lifetime = 3 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = true
	app.Session = sessionManager

	repos := handlers.NewRepo(&app)
	r := chi.NewRouter()
	r.Use(middleware.MiddleLogger)
	r.Use(middleware.NoSurf)
	r.Use(middleware.NewSessionLoad(&app))
	r.Get("/", repos.Home)
	r.Get("/about", repos.About)
	fileServer := http.FileServer(http.Dir("../assests/"))
	r.Handle("/assests/*", http.StripPrefix("/assests", fileServer))
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error while running server")
	}
}
