package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/philippta/auth"
)

type user struct {
	name string
}

var users = map[int]*user{
	1: {"alice"},
	2: {"bob"},
	3: {"charlie"},
}

func lookup(id int) (*user, bool) {
	u, ok := users[id]
	if !ok {
		return nil, false
	}
	return u, true
}

func main() {
	auth := auth.New(lookup)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ok := auth.Ok(r.Context())
		fmt.Fprintln(w, "auth.Ok:", ok)

		user, ok := auth.User(r.Context())
		fmt.Fprintln(w, "auth.User:", user, ok)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(r.Context(), 1)
	})

	mux.HandleFunc("/login-as/", func(w http.ResponseWriter, r *http.Request) {
		idstr := r.URL.Path[len("/login-as/"):]
		id, _ := strconv.Atoi(idstr)

		auth.Login(r.Context(), id)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.Logout(r.Context())
	})

	http.ListenAndServe(":8080", auth.Handler(mux))
}
