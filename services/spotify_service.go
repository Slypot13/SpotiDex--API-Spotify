package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type SimpleArtist struct {
	ID         string
	Name       string
	ImageURL   string
	Popularity int
}

type ArtistInfo struct {
	ID         string
	Name       string
	ImageURL   string
	Popularity int
	Followers  int
	Genres     []string
	SpotifyURL string
}

type AlbumInfo struct {
	ID          string
	Name        string
	ImageURL    string
	ReleaseDate string
	TotalTracks int
	SpotifyURL  string
}


type tokenResp struct {
	AccessToken string `json:"access_token"`
}

type image struct {
	URL string `json:"url"`
}

type followers struct {
	Total int `json:"total"`
}

type externalURLs struct {
	Spotify string `json:"spotify"`
}

type artistItem struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Popularity  int        `json:"popularity"`
	Images      []image    `json:"images"`
	Followers   followers  `json:"followers"`
	Genres      []string   `json:"genres"`
	ExternalURL externalURLs `json:"external_urls"`
}

type searchResp struct {
	Artists struct {
		Items []artistItem `json:"items"`
	} `json:"artists"`
}

type albumItem struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Images      []image     `json:"images"`
	ReleaseDate string      `json:"release_date"`
	TotalTracks int         `json:"total_tracks"`
	ExternalURL externalURLs `json:"external_urls"`
}

type albumsResp struct {
	Items []albumItem `json:"items"`
}


func firstImageURL(imgs []image) string {
	if len(imgs) == 0 {
		return ""
	}
	return imgs[0].URL
}

// Récupération du token
func getToken() (string, error) {
	clientID := "b2758b6fc111451ea08499f71d2ec221"
	clientSecret := "1a6559f3d05c43359f63a4ede6cd1b8e"

	if clientID == "" || clientSecret == "" {
		return "", errors.New("client id/secret manquant")
	}

	data := "grant_type=client_credentials"

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data))
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("erreur token spotify")
	}

	var tr tokenResp
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return "", err
	}

	if tr.AccessToken == "" {
		return "", errors.New("token vide")
	}

	return tr.AccessToken, nil
}

// Fonction pour appeler l'API
func spotifyGetJSON(urlStr string, target any) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("spotify api error")
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// --------- Endpoints --------- 

// ENDPOINT 1 : 
func SearchArtists(query string) ([]SimpleArtist, error) {
	// On encode la recherche pour éviter les bugs
	q := url.QueryEscape(query)
	u := "https://api.spotify.com/v1/search?q=" + q + "&type=artist&limit=12"

	var sr searchResp
	if err := spotifyGetJSON(u, &sr); err != nil {
		return nil, err
	}

	result := []SimpleArtist{}
	for _, a := range sr.Artists.Items {
		result = append(result, SimpleArtist{
			ID:         a.ID,
			Name:       a.Name,
			ImageURL:   firstImageURL(a.Images),
			Popularity: a.Popularity,
		})
	}

	return result, nil
}

// ENDPOINT 2
func GetArtist(id string) (ArtistInfo, error) {
	u := "https://api.spotify.com/v1/artists/" + url.PathEscape(id)

	var a artistItem
	if err := spotifyGetJSON(u, &a); err != nil {
		return ArtistInfo{}, err
	}

	return ArtistInfo{
		ID:         a.ID,
		Name:       a.Name,
		ImageURL:   firstImageURL(a.Images),
		Popularity: a.Popularity,
		Followers:  a.Followers.Total,
		Genres:     a.Genres,
		SpotifyURL: a.ExternalURL.Spotify,
	}, nil
}

// ENDPOINT 3 
func GetAlbums(id string) ([]AlbumInfo, error) {
	u := "https://api.spotify.com/v1/artists/" + url.PathEscape(id) + "/albums?include_groups=album,single&limit=12"

	var ar albumsResp
	if err := spotifyGetJSON(u, &ar); err != nil {
		return nil, err
	}

	result := []AlbumInfo{}
	for _, al := range ar.Items {
		result = append(result, AlbumInfo{
			ID:          al.ID,
			Name:        al.Name,
			ImageURL:    firstImageURL(al.Images),
			ReleaseDate: al.ReleaseDate,
			TotalTracks: al.TotalTracks,
			SpotifyURL:  al.ExternalURL.Spotify,
		})
	}

	return result, nil
}
