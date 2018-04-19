package internal

type Game struct {
	ID       uint64 `json:"gamePk"`
	Link     string `json:"link"`
	GameDate string `json:"gameDate"`
	Status   struct {
		StatusCode string `json:"statusCode"`
	} `json:"status"`
	Teams map[string]struct {
		Score uint `json:"score"`
		Team  struct {
			ID           uint   `json:"id"`
			Name         string `json:"name"`
			TeamName     string `json:"teamName"`
			Abbreviation string `json:"abbreviation"`
		} `json:"team"`
		Winner bool `json:"isWinner"`
	} `json:"teams"`
}

type Date struct {
	Games      []Game `json:"games"`
	DateString string `json:"date"`
}

type Root struct {
	Dates []Date `json:"dates"`
}
