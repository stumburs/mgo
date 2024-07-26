package mgo

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"os"
)

type Ngram = map[string][]string

type MarkovGenerator struct {
	Ngrams     Ngram
	SourceText string
}

// Represents how to split source text (by space/newline or by N characters)
type SplitStrategy int

const (
	SplitBySpaces SplitStrategy = iota
	SplitByNCharacters
)

// NOTE: Most functions return *MarkovGenerator to support method chaining

// Create new markov generator object
func NewMarkovGenerator() *MarkovGenerator {
	return &MarkovGenerator{
		Ngrams: make(Ngram),
	}
}

// Read text from file
func (mg *MarkovGenerator) ReadSourceFromFile(path string) *MarkovGenerator {
	// Read file
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to open file `%s`: %s\n", path, err)
	}
	mg.SourceText = string(b[:])
	return mg
}

// Construct Ngrams from source text (N - only matters when using SplitByNCharacters). SplitBySpaces | SplitByNCharacters
func (mg *MarkovGenerator) BuildNgrams(strategy SplitStrategy, N int) *MarkovGenerator {
	words := []string{}
	switch strategy {
	case SplitBySpaces:
		words = splitKeepSpaces(mg.SourceText)
	case SplitByNCharacters:
		words = splitByNCharacters(mg.SourceText, N)
	}
	for i := 0; i < len(words)-1; i++ {
		currentWord := words[i]
		nextWord := words[i+1]
		mg.Ngrams[currentWord] = append(mg.Ngrams[currentWord], nextWord)
	}
	return mg
}

// Generate text from previously constructed Ngrams
func (mg *MarkovGenerator) GenerateText(length int) string {
	currentWord := getRandomKey(mg.Ngrams)
	generatedText := currentWord
	for i := 0; i < length; i++ {
		values := mg.Ngrams[currentWord]
		if len(values) == 0 {
			break
		}
		currentWord = values[rand.Intn(len(values))]
		generatedText += currentWord
	}
	return generatedText
}

// Continue generating text after the provided string
func (mg *MarkovGenerator) Continue(text string, length int) string {

	panic("Unimplemented")
}

// Serializes to gob and writes ngrams to file
func (mg *MarkovGenerator) WriteNgrams(path string) *MarkovGenerator {
	// Serialize ngrams
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(mg.Ngrams)
	if err != nil {
		log.Fatalf("Failed to serialize ngrams: %s\n", err)
	}

	// Write gob to file
	f, err := os.Create(path)
	if err != nil {
		log.Fatal("Couldn't open file\n")
	}
	defer f.Close()

	_, err = buffer.WriteTo(f)
	if err != nil {
		log.Fatalf("Failed to write ngrams to file\n")
	}
	return mg
}

// Reads binary gob from file and de-serializes it into ngrams
func (mg *MarkovGenerator) ReadNgrams(path string) *MarkovGenerator {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Couldn't open file\n")
	}
	defer f.Close()

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(f)
	if err != nil {
		log.Fatalf("Failed to read file to buffer\n")
	}

	// Deserialize buffer into ngrams
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&mg.Ngrams)
	if err != nil {
		log.Fatalf("Failed deserializing ngrams: %s", err)
	}
	return mg
}

// Writes text to file
func (mg *MarkovGenerator) WriteText(text string, path string) *MarkovGenerator {
	err := os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		log.Printf("Failed to write output file: %s", err)
	}
	return mg
}

/*
 *	Helper functions
 */

// Returns slice of words, spaces, and new-lines created from source text
func splitKeepSpaces(s string) []string {
	result := []string{}
	word := ""
	for _, char := range s {

		if char == ' ' || char == '\n' {
			if word != "" {
				result = append(result, word)
				word = ""
			}
			result = append(result, string(char))
		} else {
			word += string(char)
		}
	}
	if word != "" {
		result = append(result, word)
	}
	return result
}

// Splits the source string into a slice of N length strings
func splitByNCharacters(s string, N int) []string {
	result := []string{}
	for i := 0; i < len(s); i += N {
		end := i + N
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

// Returns a random key in ngrams
func getRandomKey(ngram Ngram) string {
	keys := make([]string, 0, len(ngram))
	for key := range ngram {
		keys = append(keys, key)
	}
	return keys[rand.Intn(len(keys))]
}
