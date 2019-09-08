package main

import (
	"net/http"
	"./routes"
	"./utils"
)

func main() {
	utils.LoadTemplates("templates/*.html")
	router := routes.NewRouter()
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}











