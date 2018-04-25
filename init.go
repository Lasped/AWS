package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/html/*"))
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/login", login)
	r.POST("/login", login)
	r.POST("/signup", signup)
	r.GET("/signup", signup)
	r.GET("/userMain", userMain)
	r.POST("/userMain", userMain)
	r.GET("/logout", logout)

	bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	dbUsers["James"] = user{"James", bs, "test@test.com", "admin"}

	log.Fatal(http.ListenAndServe(":80", r))

}
