package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// read
func ReadJSON(filename string, store any) error {
	content, err := os.ReadFile(filepath.Join("assets", filename))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &store); err != nil {
		return err
	}
	return nil
}

// write
func WriteJSON(filename string, content any) error {
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join("assets", filename), jsonContent, 0644)
	if err != nil {
		return err
	}
	return nil
}
