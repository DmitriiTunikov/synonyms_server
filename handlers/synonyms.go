package handlers

import (
	"encoding/json"
	"net/http"
	"spread/synonymscache"
)

type synonymsHandler struct {
	SynonymsCache synonymscache.SynonymsCache
}

func NewSynonymsHandler(synonymsCache synonymscache.SynonymsCache) http.Handler {
	return &synonymsHandler{SynonymsCache: synonymsCache}
}

func (h *synonymsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	switch r.Method {

	// OPTIONS ...
	case "OPTIONS":

	// GET localhost:8080/synonyms?word=begin
	case "GET":
		word := r.URL.Query().Get("word")

		res := output{
			Word:     word,
			Synonyms: h.SynonymsCache.Get(word),
		}

		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	// POST localhost:8080/synonyms
	case "POST":
		d := json.NewDecoder(r.Body)
		i := &input{}
		err := d.Decode(i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		h.SynonymsCache.Set(i.Word, i.Synonym)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type input struct {
	Word    string `json:"word"`
	Synonym string `json:"synonym"`
}

type output struct {
	Word     string   `json:"word"`
	Synonyms []string `json:"synonyms"`
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
