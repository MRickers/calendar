package main

import (
	"calendar"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// const data

const edit_template_path = "tmpl/edit.html"
const view_template_path = "tmpl/view.html"

// global
var templates = template.Must(template.ParseFiles(view_template_path, edit_template_path))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal().Err(http.ListenAndServe(":8080", nil))
}

// Closure to validate url path before calling handlers
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath, err := getTitle(r.URL.Path)
		if err != nil {
			log.Error().Err(err).Msg("invalid path")
			http.NotFound(w, r)
			return
		}
		fn(w, r, validPath)
	}
}

func getTitle(urlPath string) (string, error) {
	m := validPath.FindStringSubmatch(urlPath)
	if m == nil {
		return "", fmt.Errorf("invalid path %s", urlPath)
	}
	return m[2], nil // The title is the second subexpression.
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	name := r.URL.Path[len("/view/"):]
	loader := calendar.LoaderFactory("file")

	user, err := loader.Load(name)

	if err != nil {
		log.Error().Err(err).Msg("load failed")
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}

	renderTemplate("view", w, &user)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	name := r.URL.Path[len("/edit/"):]
	loader := calendar.LoaderFactory("file")
	user, err := loader.Load(name)
	if err != nil {
		log.Info().Err(err).Msg("creating default user")
	}
	renderTemplate("edit", w, &user)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	name := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	loader := calendar.LoaderFactory("file")
	user, err := loader.Load(name)

	if err != nil {
		log.Info().Msg("creating default user")
		user = calendar.CreateUser(title, "default@mail.com", body)
	} else {
		log.Info().Msg("appending")
		dates := strings.Split(body, "\r\n")
		user.Dates = append(user.Dates, dates...)
	}
	err = loader.Save(&user)

	if err != nil {
		log.Error().Err(err).Msg("could not save user data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(name string, w http.ResponseWriter, user *calendar.User) {
	err := templates.ExecuteTemplate(w, name+".html", user)
	if err != nil {
		log.Error().Err(err).Msg("could not execute template")
		http.Error(w, "loading template failed", http.StatusInternalServerError)
	}
}
