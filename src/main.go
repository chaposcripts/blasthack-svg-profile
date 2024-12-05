package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/api/blasthack-profile-card", ShowUserProfileCard)
	fmt.Println("Started server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func ShowUserProfileCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}

	profileId := r.URL.Query().Get("profileId")
	layout := r.URL.Query().Get("layout")
	if len(profileId) == 0 {
		http.Error(w, "PROFILE_ID_NOT_PROVIDED", http.StatusBadRequest)
		return
	}
	if len(layout) == 0 || layout != "extended" || layout != "compact" {
		layout = "compact"
	}

	user, err := GetProfileInfo(r.URL.Query().Get("profileId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Rendering profile:", profileId)

	card, err := template.ParseFiles("../template/" + layout + "/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	err = card.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
