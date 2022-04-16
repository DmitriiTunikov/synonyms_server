package synonymscache

type set map[string]struct{}

type SynonymsCache interface {
	// Set save that words are synonyms
	Set(word, synonym string)
	// Get return all synonyms for word
	Get(word string) []string
	// Reset clean cache
	Reset()
}
