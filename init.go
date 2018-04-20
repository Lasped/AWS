package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/html/*"))
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/login", login)
	r.GET("/userMain", userMain)
	r.POST("/userMain", signup)

	log.Fatal(http.ListenAndServe(":80", r))

}
