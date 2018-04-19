package backend

import (
	"encoding/json"
	"fmt"

	"github.com/alfreddobradi/go-mlb-results/internal"
	"github.com/dgraph-io/badger"
)

type Backend struct {
	DB *badger.DB
}

func New(dir string, valueDir string) (Backend, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir

	db, err := badger.Open(opts)
	return Backend{db}, err
}

func (b *Backend) Save(root internal.Root) ([]internal.Game, error) {
	var gameJSON []byte
	var games []internal.Game
	err := b.DB.Update(func(txn *badger.Txn) error {
		var err error
		dateString := root.Dates[0].DateString
		for _, game := range root.Dates[0].Games {
			key := fmt.Sprintf("%s:%d:%d", dateString, game.Teams["away"].Team.ID, game.Teams["home"].Team.ID)
			gameJSON, err = json.Marshal(game)
			if err != nil {
				return err
			}
			txn.Set([]byte(key), gameJSON)
			games = append(games, game)
		}

		return nil
	})

	return games, err
}
