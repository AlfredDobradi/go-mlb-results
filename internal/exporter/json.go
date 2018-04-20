package exporter

import (
	"encoding/json"

	"github.com/alfreddobradi/go-mlb-results/internal"
)

// JSON is the wrapper around the specific functionality
type JSON struct {
	Options internal.Options
}

// Export returns the data JSON-formatted
func (e JSON) Export(g internal.Game) string {
	s, _ := json.Marshal(g)
	return string(s)
}
