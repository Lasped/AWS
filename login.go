package main

import (
	"net/http"

	"github.com/Lasped/go.uuid"
	"github.com/julienschmidt/httprouter"
)

type user struct {
	UserName string
	Password string
	Email    string
}

var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	HandleError(w, err)
}

func alreadyLoggedIN(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	un := dbSessions[c.Value]
	_, ok := dbUsers[un]
	return ok
}

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if alreadyLoggedIN(r) {
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}
	// process form submission
	if r.Method == http.MethodPost {

		// get form values
		un := r.FormValue("username")
		p := r.FormValue("password")
		e := r.FormValue("email")

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// Create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un

		// store user in dbUsers
		u := user{un, p, e}
		dbUsers[un] = u

		//Redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "userMain.gohtml", nil)
}
