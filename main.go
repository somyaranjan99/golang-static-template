package main

import (
	"encoding/gob"
	"fmt"
	"github/somyaranjan99/basic-go-project/cmd/web/middleware"
	"github/somyaranjan99/basic-go-project/internal/condriver"
	"github/somyaranjan99/basic-go-project/internal/helpers"
	"github/somyaranjan99/basic-go-project/internal/reservationtypes"
	"github/somyaranjan99/basic-go-project/pkg/config"
	"github/somyaranjan99/basic-go-project/pkg/handlers"
	"github/somyaranjan99/basic-go-project/pkg/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
)

var app config.AppConfig
var Infolog *log.Logger
var ErrorLog *log.Logger

func main() {
	r, err := Run()
	//	defer close(app.MailChan)

	// msg := models.MailData{
	// 	To:      "som@gmial.com",
	// 	From:    "ran@gm.com",
	// 	Subject: "Some subject",
	// 	Content: "Hello mail",
	// }
	// mailServer := sendmail.NewMAilRepo(&app)
	// mailServer.ListenForMail()

	// app.MailChan <- msg
	if err != nil {
		log.Fatal("Error while running server")
		return
	}
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error while running server")
	}
}
func Run() (http.Handler, error) {
	gob.Register(models.Reservation{})
	gob.Register(reservationtypes.User{})
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	app := config.AppConfig{}
	sessionManager := scs.New()
	sessionManager.Lifetime = 3 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = true
	app.Session = sessionManager
	helpers.NewErrorLogs(&app)
	db, err := condriver.Dbconn()
	if err != nil {
		log.Fatal(err)
	}
	repos := handlers.NewRepo(&app, db)
	// dbrepo.NewRepositoryDBHandler(&app, db)
	Infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.Infolog = Infolog
	ErrorLog = log.New(os.Stdout, "Erro\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = ErrorLog
	r := chi.NewRouter()
	r.Use(middleware.MiddleLogger)
	r.Use(middleware.NoSurf)
	r.Use(middleware.NewSessionLoad(&app))
	r.Get("/", repos.Home)
	r.Get("/about", repos.About)
	r.Get("/generals-quarters", repos.Generals)
	r.Get("/majors-suite", repos.Majors)
	r.Get("/search-availability", repos.Aavailability)
	r.Post("/search-availability", repos.PostAavailability)
	r.Post("/search-availability-query", repos.PostSearchAvailability)

	r.Get("/contact", repos.Contact)
	r.Get("/make-reservation", repos.Reservation)
	r.Post("/make-reservation", repos.PostReservation)
	r.Get("/reservation-summary", repos.ReservationSummary)
	r.Get("/choose-room/{id}", repos.ChooseRoom)
	r.Get("/user/login", repos.Login)
	r.Post("/user/login", repos.PostLogin)
	r.Get("/user/signup", repos.UserSignup)
	r.Post("/user/signup", repos.PostUserSignup)
	// r.Route("/admin",func(r chi.Router) {

	// })

	projectRoot, err := filepath.Abs(filepath.Join("..", "..")) // Goes up from cmd/web
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	staticPath := filepath.Join(projectRoot, "static")
	fmt.Println(staticPath)
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		log.Fatalf("Static directory not found at: %s", staticPath)
		return nil, err
	}

	log.Printf("Serving static files from: %s", staticPath)
	fs := http.FileServer(http.Dir(staticPath))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	return r, nil
}
