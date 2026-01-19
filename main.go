package main

import (
	"fmt"
	"net/http"

	"myspotidex/controllers"
	"myspotidex/helper"
)

func main() {
	// On charge les templates HTML
	helper.LoadTemplates()

	// Toutes les routes du site
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/collection", controllers.Collection)
	http.HandleFunc("/artist", controllers.Artist)
	http.HandleFunc("/favorites", controllers.Favorites)
	http.HandleFunc("/about", controllers.About)

	http.HandleFunc("/favorites/add", controllers.FavoritesAdd)
	http.HandleFunc("/favorites/remove", controllers.FavoritesRemove)

	fs := http.FileServer(http.Dir("../assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur lancé http://localhost:8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
