package services

import (
	"strings"
)

func SearchLocal(artists []SimpleArtist, q string) []SimpleArtist {
	q = strings.ToLower(strings.TrimSpace(q))
	// Si y'a rien on renvoie tout
	if q == "" {
		return artists
	}

	out := []SimpleArtist{}
	for _, a := range artists {
		name := strings.ToLower(a.Name)

		// On regarde si le nom contient la recherche
		if strings.Contains(name, q) {
			out = append(out, a)
			continue
		}

		// propriété 2 : genres 
	}

	return out
}

func FilterLocal(artists []SimpleArtist, genre string, minPop int, hasImage bool) []SimpleArtist {
	genre = strings.ToLower(strings.TrimSpace(genre))

	out := []SimpleArtist{}

	for _, a := range artists {
		// filtre 1 : genre 
		if genre != "" {
			// Faire le filtre plus tard si j'ai le temps
		}

		// filtre 2 : popularité
		if a.Popularity < minPop {
			continue
		}

		// filtre 3 : image
		if hasImage && a.ImageURL == "" {
			continue
		}

		out = append(out, a)
	}

	return out
}

func Paginate(items []SimpleArtist, page int, limit int) ([]SimpleArtist, int) {
	if limit <= 0 {
		limit = 30
	}
	if page <= 0 {
		page = 1
	}

	total := len(items)
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	start := (page - 1) * limit
	if start >= total {
		return []SimpleArtist{}, totalPages
	}

	end := start + limit
	if end > total {
		end = total
	}

	return items[start:end], totalPages
}
