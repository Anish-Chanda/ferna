package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anish-chanda/ferna/db"
)

func SearchSpecies(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")

		// parse limit (default 20)
		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil {
				limit = v
			}
		}

		// parse offset (default 0)
		offset := 0
		if o := r.URL.Query().Get("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil {
				offset = v
			}
		}

		species, err := database.SearchSpecies(r.Context(), q, limit, offset)
		if err != nil {
			http.Error(w, "failed to search species: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(species)
	}
}
