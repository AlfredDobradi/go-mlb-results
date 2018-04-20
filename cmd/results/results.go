package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alfreddobradi/go-mlb-results/internal"
	"github.com/alfreddobradi/go-mlb-results/internal/backend"
	"github.com/alfreddobradi/go-mlb-results/internal/exporter"
	"github.com/alfreddobradi/go-mlb-results/internal/parser"
)

func main() {
	backend, err := backend.New("/tmp/mlb", "/tmp/mlb")
	internal.Must(err, "backend: connect:")

	options := getOptions()

	games := make([]internal.Game, 0)
	games, err = backend.Results(options.Date)
	internal.Must(err, "backend: read:")

	var counter int

	if len(games) == 0 {
		root, err := parser.Parse(options.Date)
		internal.Must(err, "parser: read:")

		games, err = backend.Save(root)
		internal.Must(err, "games: save:")

		counter = 0
	} else {
		counter, err = backend.GetCounter()
		internal.Must(err, "backend: counter: get:")
	}

	if len(games) == 0 {
		fmt.Println("No games found")
		os.Exit(0)
	}

	counter = counter + 1
	next := counter % len(games)

	selected := games[next]

	var ex exporter.Plain
	ex.Options = options
	selected.Export(ex)

	err = backend.SetCounter(counter)
	internal.Must(err, "backend: counter: set:")

	backend.DB.Close()
}

func getOptions() (options internal.Options) {
	var d string
	d = os.Getenv("DATE")
	if len(d) == 0 {
		d = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	options.Date = []byte(d)

	d = strings.ToLower(os.Getenv("COLORS"))
	options.Colors = (d == "true" || d == "1")

	return
}
