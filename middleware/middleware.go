package middleware

import (
	"github/somyaranjan99/basic-go-project/pkg/config"
	"net/http"

	"github.com/justinas/nosurf"
)

// type AppConfigWrapper struct {
// 	*config.AppConfig
// }

func NewSessionLoad(app *config.AppConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if app.Session == nil {
			panic("Session is not initialized")
		}
		return app.Session.LoadAndSave(next)
	}
}

func MiddleLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}
