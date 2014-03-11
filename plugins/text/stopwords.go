package text

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

var (
	stopWords = map[string]struct{}{}
)

func stopWordCount(str string) int {
	// Prepare input
	re := regexp.MustCompile("[^\\p{Ll}\\p{Lu}\\p{Lt}\\p{Lo}\\p{Nd}\\p{Pc}\\s]")
	str = re.ReplaceAllString(str, "")
	words := strings.Fields(str)

	// Check words against stop words
	count := 0
	for _, w := range words {
		if _, ok := stopWords[strings.ToLower(w)]; ok {
			count++
			continue
		}
	}

	return count
}

func init() {
	// Load stop words
	r := bytes.NewBuffer(stopwords_en_csv())
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		stopWords[scanner.Text()] = struct{}{}
	}
}
