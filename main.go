package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/blasthack-profile-card", ShowUserProfileCard)
	http.ListenAndServe(":8080", nil)
}

func ShowUserProfileCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}

	userId := r.URL.Query().Get("profileId")
	if len(userId) == 0 {
		http.Error(w, "PROFILE_ID_NOT_PROVIDED", http.StatusBadRequest)
		return
	}

	user, err := GetProfileInfo(r.URL.Query().Get("profileId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	card, err := template.ParseFiles("index.html", "style.css")
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
