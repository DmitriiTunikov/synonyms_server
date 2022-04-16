package synonymscache

import "sync"

// CacheFastRead struct for storing synonyms
type CacheFastRead struct {
	sync.RWMutex
	wordToSynonyms map[string]set
}

// NewSynonymsCacheFastRead create CacheFastRead
func NewSynonymsCacheFastRead() SynonymsCache {
	return &CacheFastRead{wordToSynonyms: make(map[string]set)}
}

func (c *CacheFastRead) addSynonym(word, synonym string) {
	// append synonym to word synonyms
	_, ok := c.wordToSynonyms[word]
	if !ok {
		c.wordToSynonyms[word] = make(set)
	}
	c.wordToSynonyms[word][synonym] = struct{}{}

	// append synonyms of synonym to word synonyms
	_, ok = c.wordToSynonyms[synonym]
	if ok {
		for s, _ := range c.wordToSynonyms[synonym] {
			if s != word {
				c.wordToSynonyms[word][s] = struct{}{}
				// append word to synonym synonyms
				c.wordToSynonyms[s][word] = struct{}{}
			}
		}
	}
}

func (c *CacheFastRead) Set(word, synonym string) {
	if word == synonym {
		return
	}

	c.Lock()
	defer c.Unlock()
	c.addSynonym(word, synonym)
	c.addSynonym(synonym, word)
}

func (c *CacheFastRead) Get(word string) []string {
	c.RLock()
	defer c.RUnlock()

	synonyms, ok := c.wordToSynonyms[word]
	if !ok {
		return []string{}
	}

	res := make([]string, 0, len(synonyms))
	for s, _ := range synonyms {
		res = append(res, s)
	}

	return res
}

func (c *CacheFastRead) Reset() {
	c.wordToSynonyms = make(map[string]set)
}
