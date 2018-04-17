package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/html/*"))

}

func main() {
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/login", login)
	r.GET("/userMain", userMain)
	r.POST("/userMain", userMain)

	http.ListenAndServe(":80", r)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	HandleError(w, err)
}

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

	d := struct {
		Usr string
		Psw string
	}{
		Usr: user,
		Psw: pass,
	}
	tpl.ExecuteTemplate(w, "userMain.gohtml", d)
}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
