package main

import (
	"net/http"
	"time"

	"github.com/Lasped/go.uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	UserName string
	Password []byte
	Email    string
	Role     string
}

type session struct {
	un           string
	lastActivity time.Time
}

var dbUsers = map[string]user{}       // user ID, user
var dbSessions = map[string]session{} // session ID, user ID
var dbSessionsCleaned time.Time

const sessionLength int = 30

func getUser(w http.ResponseWriter, r *http.Request) user {
	//get Cookie
	c, err := r.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	// if the user exist already, get user
	var u user
	if s, ok := dbSessions[c.Value]; ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
		u = dbUsers[s.un]
	}
	return u
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if alreadyLoggedIN(w, r) {
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		// is there a UserName
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)

}

func alreadyLoggedIN(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
	}

	_, ok = dbUsers[s.un]
	// refresh session
	c.MaxAge = sessionLength

	http.SetCookie(w, c)

	return ok
}

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if alreadyLoggedIN(w, r) {
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}
	var u user

	// process form submission
	if r.Method == http.MethodPost {

		// get form values
		un := r.FormValue("username")
		p := r.FormValue("password")
		e := r.FormValue("email")
		rl := "admin"

		//username taken?
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
		dbSessions[c.Value] = session{un, time.Now()}

		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return
		}
		u := user{un, bs, e, rl}
		dbUsers[un] = u
		//Redirect
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", u)
}
func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !alreadyLoggedIN(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	c, _ := r.Cookie("session")
	//delete the session
	delete(dbSessions, c.Value)
	//remove the Cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	//clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
		dbSessionsCleaned = time.Now()
	}
}
