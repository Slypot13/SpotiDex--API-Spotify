package helper

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func LoadTemplates() {
	funcMap := template.FuncMap{
		"add": func(a int, b int) int { return a + b },
	}
	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("../templates/*.html"))
}

func Render(w http.ResponseWriter, name string, data any) {
	err := tpl.ExecuteTemplate(w, name+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
