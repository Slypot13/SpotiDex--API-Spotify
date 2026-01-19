package services

import (
	"encoding/json"
	"os"
)

type FavArtist struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

var favFile = "../data/favorites.json"

func LoadFavorites() ([]FavArtist, error) {
	b, err := os.ReadFile(favFile)
	if err != nil {
		return []FavArtist{}, err
	}

	var items []FavArtist
	if err := json.Unmarshal(b, &items); err != nil {
		return []FavArtist{}, err
	}
	return items, nil
}

func SaveFavorites(items []FavArtist) error {
	b, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(favFile, b, 0644)
}

func IsFavorite(id string) bool {
	items, err := LoadFavorites()
	if err != nil {
		return false
	}
	for _, it := range items {
		if it.ID == id {
			return true
		}
	}
	return false
}

func AddFavorite(a FavArtist) error {
	items, err := LoadFavorites()
	if err != nil {
		items = []FavArtist{}
	}

	for _, it := range items {
		if it.ID == a.ID {
			return nil 
		}
	}

	items = append(items, a)
	return SaveFavorites(items)
}

func RemoveFavorite(id string) error {
	items, err := LoadFavorites()
	if err != nil {
		return err
	}

	out := []FavArtist{}
	for _, it := range items {
		if it.ID != id {
			out = append(out, it)
		}
	}

	return SaveFavorites(out)
}
