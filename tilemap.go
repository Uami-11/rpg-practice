package main

import (
	"encoding/json"
	"os"
)

type TilemapLayerJSON struct {
	data   []int `json: "data"`
	width  int   `json: "width"`
	height int   `json: "height"`
}

type TilemapJSON struct {
	Tilemap []TilemapLayerJSON `json: "layers"`
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
