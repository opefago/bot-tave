package common

import (
	"encoding/json"
	"os"
)

func LoadTriangularPair(filename string) ([]TraingularPair, error) {
	var triangularPairs []TraingularPair

	dat, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(dat, &triangularPairs)
	return triangularPairs, nil
}
