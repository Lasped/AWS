package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Lasped/AWS/controllers"
	"github.com/Lasped/AWS/models"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	models.Tpl = template.Must(template.ParseGlob("templates/html/*"))
}

func main() {
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/login", controllers.Login)
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.GET("/signup", controllers.Signup)
	r.GET("/userMain", controllers.UserMain)
	r.POST("/userMain", controllers.UserMain)
	r.GET("/logout", controllers.Logout)

	bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	models.DbUsers["James"] = models.User{"James", bs, "test@test.com", "admin"}

	log.Fatal(http.ListenAndServe(":80", r))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := controllers.GetUser(w, r)
	err := models.Tpl.ExecuteTemplate(w, "index.gohtml", u)
	controllers.HandleError(w, err)
}
