package services

import "errors"

func GetCollection() ([]SimpleArtist, error) {
	queries := []string{"pop", "rap", "rock", "electro", "jazz"}

	all := []SimpleArtist{}
	seen := map[string]bool{}

	for _, q := range queries {
		// On cherche un peu de tout pour remplir la page
		list, err := SearchArtists(q) //endpoint 1
		if err != nil {
			return nil, errors.New("impossible de charger la collection (spotify)")
		}

		for _, a := range list {
			if a.ID == "" {
				continue
			}
			if !seen[a.ID] {
				seen[a.ID] = true
				all = append(all, a)
			}
		}
	}

	return all, nil
}
