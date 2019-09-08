package utils

import (
	"html/template"
	"net/http"
)


var templates *template.Template

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func ExecuteTemplate(res http.ResponseWriter, tmp string, data interface{}) {
	templates.ExecuteTemplate(res, tmp, data)
}