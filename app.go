package main

import (
	"calendar"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	loader := calendar.LoaderFactory("file")

	user, err := loader.Load(title)

	if err != nil {
		log.Error().Err(err).Msg("load failed")
		return
	}

	t, _ := template.ParseFiles("view.html")
	t.Execute(w, user)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	loader := calendar.LoaderFactory("file")
	user, err := loader.Load(title)
	if err != nil {
		log.Warn().Err(err).Msg("creating default user")
	}

	t, err := template.ParseFiles("edit.html")

	if err != nil {
		log.Error().Err(err).Msg("parsing template failed")
		fmt.Fprintf(w, "<h1>Internal server error</h1>")
		return
	}

	err = t.Execute(w, user)

	if err != nil {
		log.Error().Err(err).Msg("running template failed")
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	loader := calendar.LoaderFactory("file")
	user := calendar.CreateUser("Foo", "foo@mail.com", "31-12-2022-12-12-12Info: Call Bar;12-12-2022-13-13-13Info: Diner")
	loader.Save(&user)
	u, _ := loader.Load("Foo")
	fmt.Println(u)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	log.Fatal().Err(http.ListenAndServe(":8080", nil))
}
