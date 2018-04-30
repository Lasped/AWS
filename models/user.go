package models

import (
	"html/template"
	"time"
)

type User struct {
	UserName string `json: "name"`
	Password []byte `json: "password"`
	Email    string `json: "email"`
	Role     string `json: "id"`
}

type Session struct {
	Un           string
	LastActivity time.Time
}

var Tpl *template.Template

var DbUsers = map[string]User{}       // user ID, user
var DbSessions = map[string]Session{} // session ID, user ID
var DbSessionsCleaned time.Time

const SessionLength int = 30
