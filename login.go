package main

import (
"html/template"
"net/http"

"github.com/julienschmidt/httprouter"
)

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	HandleError(w, err)
}

func userMain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user := r.FormValue("user")
	pass := r.FormValue("pass")

	if
	d := struct {
		Usr string
		Psw string
	}{
		Usr: user,
		Psw: pass,
	}
	tpl.ExecuteTemplate(w, "userMain.gohtml", d)
}
