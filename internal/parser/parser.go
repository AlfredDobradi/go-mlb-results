package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alfreddobradi/go-mlb-results/internal"
)

func Parse(date []byte) (internal.Root, error) {
	var root internal.Root
	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s&hydrate=team", date)
	response, err := http.Get(url)
	if err != nil {
		return root, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return root, err
	}

	err = json.Unmarshal(body, &root)
	return root, err
}
