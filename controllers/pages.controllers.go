package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"myspotidex/helper"
	"myspotidex/services"
)

func Home(w http.ResponseWriter, r *http.Request) {
	helper.Render(w, "home", nil)
}

type CollectionData struct {
	Artists []services.SimpleArtist

	Q        string
	Genre    string
	MinPop   int
	HasImage bool

	Page       int
	Limit      int
	TotalPages int
}

func Collection(w http.ResponseWriter, r *http.Request) {
	// Je récupère les paramètres de l'URL
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	genre := strings.TrimSpace(r.URL.Query().Get("genre"))

	minPop := 0
	if v := r.URL.Query().Get("minpop"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			minPop = n
		}
	}

	hasImage := false
	if r.URL.Query().Get("hasimage") == "1" {
		hasImage = true
	}

	page := 1
	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}

	limit := 30
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && (n == 10 || n == 20 || n == 30) {
			limit = n
		}
	}

	all, err := services.GetCollection()
	if err != nil {
		helper.Render(w, "error", map[string]string{"Message": err.Error()})
		return
	}

	// FT1 recherche
	filtered := all
	if q != "" {
		filtered = services.SearchLocal(filtered, q)
	}

	// FT2 filtres 
	filtered = services.FilterLocal(filtered, genre, minPop, hasImage)

	// FT3 pagination
	pageItems, totalPages := services.Paginate(filtered, page, limit)

	helper.Render(w, "collection", CollectionData{
		Artists:    pageItems,
		Q:          q,
		Genre:      genre,
		MinPop:     minPop,
		HasImage:   hasImage,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}


type ArtistData struct {
	Artist services.ArtistInfo
	Albums []services.AlbumInfo
	IsFav  bool
}

func Artist(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	// Pas d'id, pas d'artiste
	if id == "" {
		helper.Render(w, "error", map[string]string{"Message": "Il manque l'id de l'artiste"})
		return
	}

	artist, err := services.GetArtist(id)
	if err != nil {
		helper.Render(w, "error", map[string]string{"Message": err.Error()})
		return
	}

	albums, err := services.GetAlbums(id)
	if err != nil {
		helper.Render(w, "error", map[string]string{"Message": err.Error()})
		return
	}

	isFav := services.IsFavorite(id)

	helper.Render(w, "artist", ArtistData{
		Artist: artist,
		Albums: albums,
		IsFav:  isFav,
	})
}

type FavoritesData struct {
	Items []services.FavArtist
}

func Favorites(w http.ResponseWriter, r *http.Request) {
	items, err := services.LoadFavorites()
	if err != nil {
		helper.Render(w, "error", map[string]string{"Message": err.Error()})
		return
	}
	helper.Render(w, "favorites", FavoritesData{Items: items})
}

func FavoritesAdd(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		http.Redirect(w, r, "/favorites", http.StatusSeeOther)
		return
	}

	artist, err := services.GetArtist(id)
	if err == nil {
		_ = services.AddFavorite(services.FavArtist{
			ID:       artist.ID,
			Name:     artist.Name,
			ImageURL: artist.ImageURL,
		})
	}

	back := r.Referer()
	if back == "" {
		back = "/favorites"
	}
	http.Redirect(w, r, back, http.StatusSeeOther)
}

func FavoritesRemove(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id != "" {
		_ = services.RemoveFavorite(id)
	}
	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}

func About(w http.ResponseWriter, r *http.Request) {
	helper.Render(w, "about", nil)
}
