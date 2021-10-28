package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// This script should only run when PRs are merged into main. It links the merged PR as well as linked issues
// that were closed as a result of the merge, to the latest unreleased milestone (if exists and not already linked).

type GitHubIssue struct {
	Owner string
	Repo string
	Id int
}

func (g GitHubIssue) getMilestoneId(ctx context.Context, client *github.Client, milestone string) (int, error) {
	milestones, _, err := client.Issues.ListMilestones(ctx, g.Owner, g.Repo, nil)
	if err != nil {
		return 0, fmt.Errorf("retrieving list of milestones: %+v", err)
	}

	for _, m := range milestones {
		title := *m.Title
		if title[1:] == milestone && *m.State != "closed" {
			return *m.Number, nil
		}
	}
	// TODO create milestone here?
	return 0, nil
}

func (g GitHubIssue) getLinkedIssue(ctx context.Context, client *github.Client) (int, error) {
	resp, _, _ := client.Issues.Get(ctx, g.Owner, g.Repo, g.Id)

	if resp.Body != nil {
		bodySplit := strings.Split(*resp.Body, " ")
		keywords := regexp.MustCompile(`^[fF]ix(.)?(.)?|[cC]lose(.)?|[rR]esolve(.)?`)
		issue := regexp.MustCompile(`^#[0-9]+`)

		for i, s := range bodySplit {
			if keywords.MatchString(s) {
				// check whether next element is the issue number
				next := bodySplit[i + 1]
				if issue.MatchString(next) {
					id, _ := strconv.Atoi(next[1:])
					return id, nil
				}
			}
		}
	}

	log.Printf("[DEBUG] no special keywords found in issue description")
	return 0, nil
}

func (g GitHubIssue) updateMilestone(ctx context.Context, client *github.Client, milestoneId int) error {
	issue, _, err := client.Issues.Get(ctx, g.Owner, g.Repo, g.Id)
	if err != nil {
		return fmt.Errorf("getting issue #%d: %+v", g.Id, err)
	}

	if issue.Milestone == nil && *issue.State == "closed" {
		_, _, err := client.Issues.Edit(ctx, g.Owner, g.Repo, g.Id, &github.IssueRequest{Milestone: &milestoneId})
		if err != nil {
			return fmt.Errorf("updating milestone on issue #%d: %+v", g.Id, err)
		}
		return nil
	}
	log.Printf("[DEBUG] github issue #%d already has milestone %s", g.Id, *issue.Milestone.Title)
	return nil
}

func getMilestone() (string, error) {
	f, err := os.Open("CHANGELOG.md")
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

func newGitHubClient(token string) (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), ctx
}

func run() error {
	viper.AutomaticEnv()
	token := viper.GetString("github_token")
	owner := strings.Split(viper.GetString("github_repository"), "/")[0]
	repo := strings.Split(viper.GetString("github_repository"), "/")[1]
	prId, err := strconv.Atoi(viper.GetString("pr_number"))
	if err != nil {
		return fmt.Errorf("parsing pr number: %+v", err)
	}

	pr := GitHubIssue{owner, repo, prId}
	client, ctx := newGitHubClient(token)

	milestone, err := getMilestone()
	if err != nil {
		return fmt.Errorf("getting latest milestone from CHANGELOG.md: %s", err)
	}
	if milestone == "" {
		log.Print("[DEBUG] no unreleased milestone could be found in CHANGELOG")
		return nil
	}

	milestoneId, err := pr.getMilestoneId(ctx, client, milestone)
	if err != nil {
		return fmt.Errorf("getting milestone id: %s", err)
	}
	if milestoneId == 0 {
		log.Printf("[DEBUG] no milestone for %s exists in github, or it has been closed", milestone)
		return nil
	}

	if err = pr.updateMilestone(ctx, client, milestoneId); err != nil {
		return err
	}

	liId, err := pr.getLinkedIssue(ctx, client)
	if err != nil {
		return fmt.Errorf("getting linked issues for #%d: %+v", pr.Id, err)
	}

	if liId != 0 {
		li := GitHubIssue{owner, repo, liId}
		if err = li.updateMilestone(ctx, client, milestoneId); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}