package subreddit

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var ErrSubredditNotFound = fmt.Errorf("subreddit not found")

var subRedditList map[string]bool

func LoadSubreddits(file string) error {
	subRedditList = make(map[string]bool)

	subredditsFile, err := os.Open(file)
	if err != nil {
		log.Errorf("error opening subreddits file: %s", err)
		return fmt.Errorf("error opening subreddits file: %s", err)
	}
	defer subredditsFile.Close()

	fileScanner := bufio.NewScanner(subredditsFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		if !strings.Contains(fileScanner.Text(), "#") {
			subRedditList[strings.ToLower(strings.TrimSpace(fileScanner.Text()))] = true
		}
	}

	log.Printf("loaded %d subreddits from file", len(subRedditList))

	return nil
}

func Exists(subreddit string) bool {
	return subRedditList[strings.ToLower(strings.TrimSpace(subreddit))]
}

func Set(subreddit string) {
	if subRedditList == nil {
		subRedditList = make(map[string]bool)
	}

	subRedditList[strings.ToLower(strings.TrimSpace(subreddit))] = true
}
