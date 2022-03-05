package main

import (
	"calendar"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	fmt.Fprintf(w, "<h1>%s</h1><div>E-Mail: %s</div><div>%s</div>", user.Name, user.Mail, strings.Join(user.Dates, "\r\n"))
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	http.HandleFunc("/view/", viewHandler)
	log.Fatal().Err(http.ListenAndServe(":8080", nil))
}
