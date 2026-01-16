package loader

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// WorktreeLoader manages a temporary git worktree for PR analysis
type WorktreeLoader struct {
	prNumber     int
	owner        string
	repo         string
	remoteName   string
	worktreePath string
	baseBranch   string
}

// NewWorktreeLoader creates a new WorktreeLoader
func NewWorktreeLoader(prNumber int, remoteName, baseBranch string) *WorktreeLoader {
	owner, repo := getRepoInfo()
	return &WorktreeLoader{
		prNumber:   prNumber,
		remoteName: remoteName,
		baseBranch: baseBranch,
		owner:      owner,
		repo:       repo,
	}
}

// detectRemoteForPR detects the appropriate remote for fetching PR
// Priority: user-specified remote flag > upstream > origin
func (l *WorktreeLoader) detectRemoteForPR() (string, error) {
	// Check if user specified remote
	if l.remoteName != "" {
		log.Printf("Using user-specified remote: %s", l.remoteName)
		return l.remoteName, nil
	}

	// Auto-detect: prefer upstream, fallback to origin
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("failed to list remotes: %w", err)
	}

	var foundOrigin, foundUpstream bool
	for _, remote := range remotes {
		name := remote.Config().Name
		if name == "upstream" {
			foundUpstream = true
		}
		if name == "origin" {
			foundOrigin = true
		}
	}

	// Prefer upstream (main repo) over origin (could be fork)
	if foundUpstream {
		log.Printf("Auto-detected remote: upstream")
		return "upstream", nil
	}

	if foundOrigin {
		log.Printf("Auto-detected remote: origin")
		return "origin", nil
	}

	return "", fmt.Errorf("no suitable remote found (origin or upstream)")
}

// Setup fetches the PR and creates a temporary worktree
func (l *WorktreeLoader) Setup() (string, error) {
	// 0. Verify we're in a git repository
	if _, err := git.PlainOpen("."); err != nil {
		return "", fmt.Errorf("not in a git repository. Please run this tool from the terraform-provider-azurerm directory")
	}

	// 1. Detect which remote to use (if not specified)
	if l.remoteName == "" {
		remote, err := l.detectRemoteForPR()
		if err != nil {
			return "", fmt.Errorf("failed to detect remote: %w. Please ensure you're in the terraform-provider-azurerm repository with correct remote configured", err)
		}
		l.remoteName = remote
	}

	// 2. Get PR details to find the base branch (if not specified)
	if l.baseBranch == "" {
		if err := l.fetchPRDetails(); err != nil {
			return "", fmt.Errorf("failed to fetch PR details: %w", err)
		}
	}

	// 3. Fetch the PR ref
	prRef := fmt.Sprintf("refs/pull/%d/head", l.prNumber)
	log.Printf("Fetching PR #%d from remote '%s' (%s/%s)...", l.prNumber, l.remoteName, l.owner, l.repo)

	cmd := exec.Command("git", "fetch", "--depth=1", l.remoteName, prRef)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to fetch PR ref: %w\n%s", err, string(output))
	}

	// 4. Create temporary worktree directory
	l.worktreePath = filepath.Join(os.TempDir(), fmt.Sprintf("azurerm-linter-pr-%d", l.prNumber))

	// Clean up any existing worktree with the same name
	if _, err := os.Stat(l.worktreePath); err == nil {
		err = l.Cleanup()
		if err != nil {
			return "", err
		}
	}

	// 5. Create the worktree
	cmd = exec.Command("git", "worktree", "add", l.worktreePath, "FETCH_HEAD")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create worktree: %w\n%s", err, string(output))
	}

	log.Printf("✓ Worktree created at %s", l.worktreePath)
	return l.worktreePath, nil
}

// fetchPRDetails fetches PR information from GitHub API to get base branch
func (l *WorktreeLoader) fetchPRDetails() error {
	prInfo, err := fetchPRInfo(os.Getenv("GITHUB_TOKEN"), l.owner, l.repo, l.prNumber)
	if err != nil {
		return err
	}

	l.baseBranch = prInfo.Base.Ref

	return nil
}

// GetWorktreePath returns the path to the worktree
func (l *WorktreeLoader) GetWorktreePath() string {
	return l.worktreePath
}

// GetBaseBranch returns the base branch of the PR
func (l *WorktreeLoader) GetBaseBranch() string {
	return l.baseBranch
}

// Cleanup removes the temporary worktree
func (l *WorktreeLoader) Cleanup() error {
	if l.worktreePath == "" {
		return nil
	}

	log.Printf("Cleaning up worktree at %s...", l.worktreePath)

	// First try to remove the worktree using git
	cmd := exec.Command("git", "worktree", "remove", l.worktreePath, "--force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: git worktree remove failed: %v\n%s", err, string(output))

		// Fallback: try to remove the directory directly
		if removeErr := os.RemoveAll(l.worktreePath); removeErr != nil {
			return fmt.Errorf("failed to remove worktree directory: %w", removeErr)
		}

		// Also prune the worktree from git's records
		cmd = exec.Command("git", "worktree", "prune")
		if pruneErr := cmd.Run(); pruneErr != nil {
			log.Printf("Warning: failed to prune worktrees: %v", pruneErr)
		}
	}

	log.Println("✓ Worktree cleanup complete")
	return nil
}
