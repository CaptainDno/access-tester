package scraper

import (
	"access-tester/src/common"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Scraper struct {
	phrases          []string
	filterTags       []*regexp.Regexp
	maxContentLength int64
}

func loadFilteredTagsArray() []*regexp.Regexp {
	filterTagsRaw, ok := os.LookupEnv("FILTER_TAGS")
	if !ok {
		filterTagsRaw = "script"
	}
	filterTags := strings.Split(filterTagsRaw, ",")

	filterTagsRegexps := make([]*regexp.Regexp, len(filterTags))
	for index, tag := range filterTags {
		filterTagsRegexps[index] = regexp.MustCompile(fmt.Sprintf("<%s.*?>(.*)</%s>", tag, tag))
	}
	return filterTagsRegexps
}

func NewScraper() *Scraper {
	maxContentLengthRaw, ok := os.LookupEnv("MAX_CONTENT_LENGTH")
	var maxContentLength int64
	if !ok {
		maxContentLength = 5_000_000
	} else {
		var err error
		maxContentLength, err = strconv.ParseInt(maxContentLengthRaw, 10, 64)
		if err != nil {
			maxContentLength = 5_000_000
		}
	}
	log.Printf("[scraper] initializing with MAX_CONTENT_LENGTH=%d", maxContentLength)
	filePath, ok := os.LookupEnv("PHRASES_FILE")
	if !ok {
		filePath = "default_phrases.txt"
	}
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var phrases []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// normalize to lowercase
		phrases = append(phrases, strings.ToLower(scanner.Text()))
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	log.Printf("[scraper] loaded %d trigger phrases from file %s", len(phrases), filePath)
	return &Scraper{
		phrases:          phrases,
		maxContentLength: maxContentLength,
		filterTags:       loadFilteredTagsArray(),
	}

}

func (s *Scraper) CheckIfBlocked(url string) (common.Result, error) {
	result := common.Result{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	req.Header.Set("Accept", "text/html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	result.Status = res.StatusCode
	result.CanBeBlocked = res.StatusCode == 403

	if res.ContentLength > s.maxContentLength {
		return result, errors.New("too big response")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	// Prepare body
	bodyString := strings.ToLower(string(bodyBytes))
	for _, expr := range s.filterTags {
		bodyString = expr.ReplaceAllString(bodyString, "")
	}

	// Check for every phrase
	phrasesFound := make([]string, 0)
	for _, str := range s.phrases {
		if strings.Contains(bodyString, str) {
			phrasesFound = append(phrasesFound, str)
		}
	}
	result.TriggerPhrasesFound = phrasesFound
	if len(phrasesFound) > 0 {
		result.CanBeBlocked = true
	}

	return result, nil
}
