package main

import (
	"net/http"
	"spread/handlers"
	"spread/synonymscache"
)

func main() {
	http.Handle("/synonyms", handlers.NewSynonymsHandler(synonymscache.NewSynonymsCacheFastRead()))
	http.ListenAndServe(":8081", nil)
}
