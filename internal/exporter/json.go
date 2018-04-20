package exporter

import (
	"encoding/json"

	"github.com/alfreddobradi/go-mlb-results/internal"
)

type JSON struct {
	Options internal.Options
}

func (e JSON) Export(g internal.Game) string {
	s, _ := json.Marshal(g)
	return string(s)
}
