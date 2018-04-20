package exporter

import (
	"fmt"

	"github.com/alfreddobradi/go-mlb-results/internal"
)

type Plain struct {
	Options internal.Options
}

func (e Plain) Export(g internal.Game) string {
	home := g.Teams["home"]
	away := g.Teams["away"]

	homeStr := fmt.Sprintf("%d %s", home.Score, home.Team.TeamName)
	awayStr := fmt.Sprintf("%s %d", away.Team.TeamName, away.Score)
	if e.Options.Colors && home.Winner {
		homeStr = fmt.Sprintf("<span color='green'>%s</span>", homeStr)
	} else if e.Options.Colors && away.Winner {
		awayStr = fmt.Sprintf("<span color='green'>%s</span>", awayStr)
	}

	return fmt.Sprintf(
		"%s : %s",
		awayStr,
		homeStr,
	)
}
