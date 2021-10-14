package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
)

// This go Script performs the following actions
// 1. Get the current milestone
// 2. Link the milestone on PRs merged into main
// 3. Link the milestone on issues linked to the PR
// 4. Go through the CHANGELOG and retrospectively perform 2. and 3.
// This script will be called by make
// This script will be triggered by github actions

func getMilestone() (string, error) {
	f, err := os.Open("../../../CHANGELOG.md")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	milestone := ""

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Unreleased") {
			for _, s := range strings.Split(scanner.Text(), " ") {
				r := regexp.MustCompile(`^[0-9].[0-9]+.[0-9]`)
				if r.MatchString(s) {
					milestone = s
				}
			}
			break
		}
	}

	return milestone, nil
}

func main() {
	pr := viper.GetString("")

	milestone, err := getMilestone()
	if err != nil {
		log.Printf("failed to link milestone: %s", err)
		os.Exit(1)
	}

	log.Printf("[DEBUG] milestone is %s", milestone)

	os.Exit(0)
}