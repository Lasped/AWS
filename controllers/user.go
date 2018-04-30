package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/Lasped/AWS/models"
	uuid "github.com/Lasped/go.uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func UserMain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := GetUser(w, r)

	if AlreadyLoggedIN(w, r) {
		err := models.Tpl.ExecuteTemplate(w, "userMain.gohtml", u)
		HandleError(w, err)
		return
	} else {
		err := models.Tpl.ExecuteTemplate(w, "login.gohtml", u)
		HandleError(w, err)
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) models.User {
	//get Cookie
	c, err := r.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	c.MaxAge = models.SessionLength
	http.SetCookie(w, c)

	// if the user exist already, get user
	var u models.User
	if s, ok := models.DbSessions[c.Value]; ok {
		s.LastActivity = time.Now()
		models.DbSessions[c.Value] = s
		u = models.DbUsers[s.Un]
	}
	return u
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if AlreadyLoggedIN(w, r) {
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		// is there a UserName
		u, ok := models.DbUsers[un]
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
		models.DbSessions[c.Value] = models.Session{un, time.Now()}
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}

	models.Tpl.ExecuteTemplate(w, "login.gohtml", nil)

}

func AlreadyLoggedIN(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	s, ok := models.DbSessions[c.Value]
	if ok {
		s.LastActivity = time.Now()
		models.DbSessions[c.Value] = s
	}

	_, ok = models.DbUsers[s.Un]
	// refresh session
	c.MaxAge = models.SessionLength

	http.SetCookie(w, c)

	return ok
}

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if AlreadyLoggedIN(w, r) {
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}
	var u models.User

	// process form submission
	if r.Method == http.MethodPost {

		// get form values
		un := r.FormValue("username")
		p := r.FormValue("password")
		e := r.FormValue("email")
		rl := "admin"

		//username taken?
		if _, ok := models.DbUsers[un]; ok {
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
		models.DbSessions[c.Value] = models.Session{un, time.Now()}

		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return
		}
		u := models.User{un, bs, e, rl}
		models.DbUsers[un] = u
		//Redirect
		http.Redirect(w, r, "/userMain", http.StatusSeeOther)
		return
	}
	models.Tpl.ExecuteTemplate(w, "signup.gohtml", u)
}
func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !AlreadyLoggedIN(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	c, _ := r.Cookie("session")
	//delete the session
	delete(models.DbSessions, c.Value)
	//remove the Cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	//clean up dbSessions
	if time.Now().Sub(models.DbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func cleanSessions() {
	for k, v := range models.DbSessions {
		if time.Now().Sub(v.LastActivity) > (time.Second * 30) {
			delete(models.DbSessions, k)
		}
		models.DbSessionsCleaned = time.Now()
	}
}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
