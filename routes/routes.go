 package routes

 import (
 	"net/http"
 	"github.com/gorilla/mux"
 	"../models"
 	"../sessions"
 	"../utils"
 	"../helpers"
 )

func NewRouter() *mux.Router {
 	router := mux.NewRouter()
	router.HandleFunc("/", helpers.AuthRequired(indexGetHandler)).Methods("GET")
	router.HandleFunc("/", helpers.AuthRequired(indexPostHandler)).Methods("POST")
	router.HandleFunc("/login", loginGetHandler).Methods("Get")
	router.HandleFunc("/login", loginPostHandler).Methods("POST")
	router.HandleFunc("/register", registerGetHandler).Methods("GET")
	router.HandleFunc("/register", registerPostHandler).Methods("POST")
	fileServer := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	return router
}


func indexGetHandler(res http.ResponseWriter, req *http.Request) {
	utils.ExecuteTemplate(res, "index.html", nil)
}

func indexPostHandler(res http.ResponseWriter, req *http.Request) {
}

func loginGetHandler(res http.ResponseWriter, req *http.Request) {
	utils.ExecuteTemplate(res, "login.html", nil)
}

func loginPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(res, "login.html", nil)
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(res, "login.html", nil)
		default:
			res.Write([]byte("Server Error"))
			return	
		}
	}


	
	session, _ := sessions.Store.Get(req, "session") 
	session.Values["username"] = username
	session.Save(req, res)
	http.Redirect(res, req, "/", 302)
	
}
func registerGetHandler(res http.ResponseWriter, req *http.Request) {
	utils.ExecuteTemplate(res, "register.html", nil)
}

func registerPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(res, req, "/login", 302)
}