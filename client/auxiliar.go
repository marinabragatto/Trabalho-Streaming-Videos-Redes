package client

import (
	"encoding/json"
	// "os"
)

type Manifest struct {
	Segments []string `json:"segments"`
}

func ReadManifest(buffer []byte) ([]string, error){

	var manifest Manifest

	err := json.Unmarshal(buffer, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest.Segments, nil
}
