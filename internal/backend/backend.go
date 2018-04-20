package backend

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alfreddobradi/go-mlb-results/internal"
	"github.com/dgraph-io/badger"
)

// Backend struct holds the DB connection
type Backend struct {
	DB *badger.DB
}

// New returns a new badger connection
func New(dir string, valueDir string) (Backend, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir

	db, err := badger.Open(opts)
	return Backend{db}, err
}

// Save saves a collection of games
func (b *Backend) Save(root internal.Root) ([]internal.Game, error) {
	var gameJSON []byte
	var games []internal.Game
	err := b.DB.Update(func(txn *badger.Txn) error {
		var err error
		for _, date := range root.Dates {
			dateString := date.DateString
			for _, game := range date.Games {
				key := fmt.Sprintf("%s:%d:%d", dateString, game.Teams["away"].Team.ID, game.Teams["home"].Team.ID)
				gameJSON, err = json.Marshal(game)
				if err != nil {
					return err
				}
				txn.Set([]byte(key), gameJSON)
				games = append(games, game)
			}
		}

		return nil
	})

	return games, err
}

// Results returns a collection of results based on the prefix
func (b *Backend) Results(prefix []byte) ([]internal.Game, error) {
	var games []internal.Game

	err := b.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}

			var g internal.Game
			err = json.Unmarshal(v, &g)
			if err != nil {
				return err
			}

			games = append(games, g)
		}

		it.Close()

		return nil
	})

	return games, err
}

// GetCounter returns the current counter - used for rotating results
func (b *Backend) GetCounter() (counter int, err error) {
	err = b.DB.View(func(txn *badger.Txn) error {
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

	return
}

// SetCounter sets a new counter
func (b *Backend) SetCounter(val int) (err error) {
	c := []byte(strconv.Itoa(val))
	b.DB.Update(func(txn *badger.Txn) error {
		txn.Set([]byte("counter"), c)
		return nil
	})

	return
}
