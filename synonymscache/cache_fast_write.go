package synonymscache

import "sync"

// CacheFastWrite struct for storing synonyms
type CacheFastWrite struct {
	sync.RWMutex
	wordToSynonyms map[string]set
}

// NewSynonymsCacheFastWrite create CacheFastWrite
func NewSynonymsCacheFastWrite() SynonymsCache {
	return &CacheFastWrite{wordToSynonyms: make(map[string]set)}
}

func (c *CacheFastWrite) addSynonyms(word, synonym string) {
	_, ok := c.wordToSynonyms[word]
	if !ok {
		c.wordToSynonyms[word] = make(set)
	}
	c.wordToSynonyms[word][synonym] = struct{}{}
}

func (c *CacheFastWrite) Set(word, synonym string) {
	c.Lock()
	defer c.Unlock()

	c.addSynonyms(word, synonym)
	c.addSynonyms(synonym, word)
}

func (c *CacheFastWrite) Get(word string) []string {
	c.RLock()
	defer c.RUnlock()
	// if we do not have synonyms for word -> just return empty slice
	_, ok := c.wordToSynonyms[word]
	if !ok {
		return []string{}
	}

	// visited words
	visited := make(map[string]struct{})
	// put word on stack
	stack := []string{word}
	for len(stack) != 0 {
		// pop word from stack
		elem := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// save it to visited map
		visited[elem] = struct{}{}
		// push all not visited word synonyms on stack
		synonyms := c.wordToSynonyms[elem]
		for synonym, _ := range synonyms {
			_, ok = visited[synonym]
			if !ok {
				stack = append(stack, synonym)
			}
		}
	}

	res := make([]string, 0, len(visited)-1)
	for w, _ := range visited {
		if w != word {
			res = append(res, w)
		}
	}

	return res
}

func (c *CacheFastWrite) Reset() {
	c.wordToSynonyms = make(map[string]set)
}
