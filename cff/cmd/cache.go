/*
Copyright © 2026 Aurélien Bulliard
*/
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)


type TripCache struct {
	Timestamp   time.Time    `json:"timestamp"`
	From        string       `json:"from"`
	To          string       `json:"to"`
	Connections []Connection `json:"connections"`
}

//Gets the full cache file path
func cacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "cff-cli", "last_trip.json"), nil
}


//Saves the data to the cache
func saveTripCache(from, to string, connections []Connection) error {
	path, err := cacheFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	cache := TripCache{
		Timestamp:   time.Now(),
		From:        from,
		To:          to,
		Connections: connections,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

//Loads cache json
func loadTripCache() (TripCache, error) {
	path, err := cacheFilePath()
	if err != nil {
		return TripCache{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return TripCache{}, err
	}

	var cache TripCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return TripCache{}, err
	}

	return cache, nil
}