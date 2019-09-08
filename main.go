package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

)

var client *redis.Client
var templates *template.Template
var store = sessions.NewCookieStore([]byte("n1XZu33ku4V"))

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	router := mux.NewRouter()
	router.HandleFunc("/", AuthRequired(indexGetHandler)).Methods("GET")
	router.HandleFunc("/", AuthRequired(indexPostHandler)).Methods("POST")
	router.HandleFunc("/login", loginGetHandler).Methods("Get")
	router.HandleFunc("/login", loginPostHandler).Methods("POST")
	router.HandleFunc("/register", registerGetHandler).Methods("GET")
	router.HandleFunc("/register", registerPostHandler).Methods("POST")
	fileServer := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)

}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request){
		session, _ := store.Get(req, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(res, req, "/login", 302)
			return
		}
		handler.ServeHTTP(res, req)
	}
}

func indexGetHandler(res http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(res, "index.html", nil)
}

func indexPostHandler(res http.ResponseWriter, req *http.Request) {
}

func loginGetHandler(res http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(res, "login.html", nil)
}

func loginPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil {
		templates.ExecuteTemplate(res, "login.html", "Unknown User")
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal Server Error"))
		return
	}	
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		templates.ExecuteTemplate(res, "login.html", "Invalid Login")
		return
	}
	session, _ := store.Get(req, "session") 
	session.Values["username"] = username
	session.Save(req, res)
	http.Redirect(res, req, "/", 302)
	
}
func registerGetHandler(res http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(res, "register.html", nil)
}

func registerPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal Server Error"))
		return
	}
	err = client.Set("user:" + username, hash, 0).Err()
		if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(res, req, "/login", 302)
}


// func testGetHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _:= store.Get(r, "session")
// 	untyped, ok := session.Valuse["username"]
// 	if ok {
// 		return
// 	}
// }










