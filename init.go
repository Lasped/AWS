package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/html/*"))
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/login", login)
	r.GET("/userMain", userMain)
	r.POST("/userMain", userMain)

	http.ListenAndServe(":80", r)
}
