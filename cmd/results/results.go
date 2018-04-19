package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/alfreddobradi/go-mlb-results/internal"
	"github.com/alfreddobradi/go-mlb-results/internal/backend"
	"github.com/alfreddobradi/go-mlb-results/internal/parser"
	"github.com/dgraph-io/badger"
)

func main() {

	backend, err := backend.New("/tmp/mlb", "/tmp/mlb")
	if err != nil {
		log.Printf("error: backend: %+v", err)
		os.Exit(1)
	}

	now := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	games := make([]internal.Game, 0)
	var counter int
	err = backend.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(now)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}

			var g internal.Game
			err = json.Unmarshal(v, &g)

			games = append(games, g)
		}

		counterItem, err := txn.Get([]byte("counter"))
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}

			counter = 0
			return nil
		}

		counterValue, _ := counterItem.Value()
		counter, _ = strconv.Atoi(string(counterValue))

		return nil
	})

	if len(games) == 0 {
		root, err := parser.Parse(now)

		if err != nil {
			log.Printf("games: save: %+v", err)
			os.Exit(1)
		}

		games, err = backend.Save(root)

		if err != nil {
			log.Printf("games: save: %+v", err)
			os.Exit(1)
		}

		counter = 0
	}

	if len(games) == 0 {
		fmt.Println("No games found")
		os.Exit(0)
	}

	counter = counter + 1
	next := counter % len(games)

	selected := games[next]
	fmt.Printf("%s %d : %d %s\n", selected.Teams["away"].Team.TeamName, selected.Teams["away"].Score, selected.Teams["home"].Score, selected.Teams["home"].Team.TeamName)

	backend.DB.Update(func(txn *badger.Txn) error {
		txn.Set([]byte("counter"), []byte(strconv.Itoa(counter)))
		return nil
	})

	backend.DB.Close()

}
