package main

import (
	"net/http"
	"chit_chat/data"
)

func login(w http.ResponseWriter, r *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(w, nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

func signupAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		danger(err, "cannot parse form")
	}

	user := data.User{
		Name: r.PostFormValue("name"),
		Email: r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	if err := user.Create(); err != nil {
		danger(err, "cannot create user")
	}

	http.Redirect(w, r, "/login", 302)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		danger(err, "cannot find user")
	}

	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "cannot create session")
		}
		cookie := http.Cookie{
			Name: "_cookie",
			Value: session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
