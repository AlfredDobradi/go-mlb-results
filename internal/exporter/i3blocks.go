package exporter

import (
	"fmt"

	"github.com/alfreddobradi/go-mlb-results/internal"
)

type I3Blocks struct {
	Options internal.Options
}

func (e I3Blocks) Export(g internal.Game) string {
	home := g.Teams["home"]
	away := g.Teams["away"]
	homeLongStr := fmt.Sprintf("%d %s", home.Score, home.Team.Name)
	awayLongStr := fmt.Sprintf("%s %d", away.Team.Name, away.Score)
	homeStr := fmt.Sprintf("%d %s", home.Score, home.Team.TeamName)
	awayStr := fmt.Sprintf("%s %d", away.Team.TeamName, away.Score)
	if e.Options.Colors && home.Winner {
		homeLongStr = fmt.Sprintf("<span color='green'>%s</span>", homeLongStr)
		homeStr = fmt.Sprintf("<span color='green'>%s</span>", homeStr)
	} else if e.Options.Colors && away.Winner {
		awayLongStr = fmt.Sprintf("<span color='green'>%s</span>", awayLongStr)
		awayStr = fmt.Sprintf("<span color='green'>%s</span>", awayStr)
	}
	long := fmt.Sprintf("%s : %s", awayLongStr, homeLongStr)
	short := fmt.Sprintf("%s : %s", awayStr, homeStr)

	return fmt.Sprintf(
		"%s\n%s",
		long,
		short,
	)
}
