package synonymscache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type synonyms struct {
	word    string
	synonym string
}

func TestSynonymsCache(t *testing.T) {
	tests := []struct {
		name             string
		initialSynonyms  []synonyms
		word             string
		expectedSynonyms []string
	}{
		{
			name: "simple case",
			initialSynonyms: []synonyms{
				{
					word:    "begin",
					synonym: "start",
				},
				{
					word:    "begin",
					synonym: "initiate",
				},
			},
			word:             "begin",
			expectedSynonyms: []string{"start", "initiate"},
		},
		{
			name: "transitive case",
			initialSynonyms: []synonyms{
				{
					word:    "a",
					synonym: "b",
				},
				{
					word:    "b",
					synonym: "c",
				},
				{
					word:    "c",
					synonym: "d",
				},
				{
					word:    "c",
					synonym: "e",
				},
			},
			word:             "a",
			expectedSynonyms: []string{"b", "c", "d", "e"},
		},
		{
			name: "synonyms not found",
			initialSynonyms: []synonyms{
				{
					word:    "begin",
					synonym: "start",
				},
				{
					word:    "begin",
					synonym: "initiate",
				},
			},
			word:             "okey",
			expectedSynonyms: []string{},
		},
		{
			name: "same synonyms",
			initialSynonyms: []synonyms{
				{
					word:    "a",
					synonym: "b",
				},
				{
					word:    "a",
					synonym: "b",
				},
			},
			word:             "a",
			expectedSynonyms: []string{"b"},
		},
		{
			name: "circle synonyms",
			initialSynonyms: []synonyms{
				{
					word:    "a",
					synonym: "b",
				},
				{
					word:    "b",
					synonym: "a",
				},
			},
			word:             "a",
			expectedSynonyms: []string{"b"},
		},
		{
			name: "word->word synonym",
			initialSynonyms: []synonyms{
				{
					word:    "a",
					synonym: "a",
				},
				{
					word:    "a",
					synonym: "b",
				},
			},
			word:             "a",
			expectedSynonyms: []string{"b"},
		},
	}

	cacheImpls := map[string]SynonymsCache{
		"fast_read_cache":  NewSynonymsCacheFastRead(),
		"fast_write_cache": NewSynonymsCacheFastWrite(),
	}
	for _, tt := range tests {
		for cacheImplName, cache := range cacheImpls {
			t.Run(fmt.Sprintf("%s(%s)", tt.name, cacheImplName), func(t *testing.T) {
				cache.Reset()
				for _, s := range tt.initialSynonyms {
					cache.Set(s.word, s.synonym)
				}

				res := cache.Get(tt.word)
				assert.ElementsMatch(t, tt.expectedSynonyms, res)
			})
		}
	}
}
