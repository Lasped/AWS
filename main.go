package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func main() {

}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := getUser(w, r)
	err := tpl.ExecuteTemplate(w, "index.gohtml", u)
	HandleError(w, err)
}

func userMain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := getUser(w, r)

	if !alreadyLoggedIN(r) {
		err := tpl.ExecuteTemplate(w, "index.gohtml", u)
		HandleError(w, err)
		return
	}
	tpl.ExecuteTemplate(w, "userMain.gohtml", u)

}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
