package helpers

import (
	"net/http"
	"../sessions"
)



func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request){
		session, _ := sessions.Store.Get(req, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(res, req, "/login", 302)
			return
		}
		handler.ServeHTTP(res, req)
	}
}