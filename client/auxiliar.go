package client

import (
	"encoding/json"
	"os"
)

type Manifest struct {
	Segments []string `json:"segments"`
}

func ReadManifest() ([]string, error) {

	file, err := os.Open("./client/segments/manifest.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var manifest Manifest

	err = json.NewDecoder(file).Decode(&manifest)
	if err != nil {
		return nil, err
	}

	return manifest.Segments, nil
}
