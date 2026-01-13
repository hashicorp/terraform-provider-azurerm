package loader

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// LocalGitLoader loads changes from local git repository
type LocalGitLoader struct {
	remoteName string
	baseBranch string
}

// Load loads changes from local git repository and returns a ChangeSet
func (l *LocalGitLoader) Load() (*ChangeSet, error) {
	cs := NewChangeSet()

	// Find git repository root
	gitRoot, err := FindGitRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find git repository: %w", err)
	}

	repo, err := git.PlainOpen(gitRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	targetCommit, _, err := resolveForLocal(repo, l.remoteName, l.baseBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve target: %w", err)
	}

	if err := processDiffWithWorktree(cs, targetCommit, gitRoot); err != nil {
		return nil, fmt.Errorf("failed to parse diff: %w", err)
	}

	log.Printf("✓ Found %d changed files with %d changed lines",
		len(cs.changedFiles), cs.getTotalChangedLines())

	return cs, nil
}

// processDiffWithWorktree compares a commit with the current worktree using git diff
func processDiffWithWorktree(cs *ChangeSet, baseCommit *object.Commit, gitRoot string) error {
	cmd := exec.Command("git", "diff", baseCommit.Hash.String())
	cmd.Dir = gitRoot // Set working directory to git root
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run git diff: %w", err)
	}

	diffOutput := string(output)
	if diffOutput == "" {
		return nil
	}

	return cs.parseDiffOutput(diffOutput)
}

// resolveForLocal resolves the target commit and worktree for comparison
func resolveForLocal(repo *git.Repository, remoteName, baseBranch string) (*object.Commit, *git.Worktree, error) {
	head, err := repo.Head()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	if !head.Name().IsBranch() {
		return nil, nil, fmt.Errorf("not on a branch (detached HEAD)")
	}

	currentBranch := head.Name().Short()
	log.Printf("Current branch: %s", currentBranch)

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get head commit: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	targetRemote, targetBranch, err := detectTargetBranch(repo, currentBranch, remoteName, baseBranch)
	if err != nil {
		return nil, nil, err
	}

	targetRef, err := repo.Reference(
		plumbing.NewRemoteReferenceName(targetRemote, targetBranch),
		true,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get %s/%s: %w", targetRemote, targetBranch, err)
	}

	targetCommit, err := repo.CommitObject(targetRef.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get target commit: %w", err)
	}

	mergeBases, err := headCommit.MergeBase(targetCommit)
	if err != nil || len(mergeBases) == 0 {
		log.Printf("Warning: failed to find merge-base, using target directly: %v", err)
		return targetCommit, worktree, nil
	}

	mergeBase := mergeBases[0]
	log.Printf("Merge-base: %s", mergeBase.Hash.String()[:7])

	return mergeBase, worktree, nil
}

// detectTargetBranch detects the target remote and branch for comparison
func detectTargetBranch(repo *git.Repository, currentBranch, remoteName, baseBranch string) (string, string, error) {
	var detectedRemote, detectedBranch string

	// Check user-specified options
	if remoteName != "" {
		detectedRemote = remoteName
	}
	if baseBranch != "" {
		detectedBranch = baseBranch
	}

	if detectedRemote == "" || detectedBranch == "" {
		if configRemote, configBranch, ok := getUpstreamFromConfig(repo, currentBranch); ok {
			if detectedRemote == "" {
				detectedRemote = configRemote
			}
			if detectedBranch == "" {
				detectedBranch = configBranch
			}
			if detectedRemote == configRemote && detectedBranch == configBranch {
				log.Printf("Using upstream from branch config: %s/%s", detectedRemote, detectedBranch)
			}
		}
	}

	if detectedRemote == "" {
		remote, err := autoDetectRemote(repo)
		if err != nil {
			return "", "", err
		}
		detectedRemote = remote
	}

	if detectedBranch == "" {
		detectedBranch = "main"
	}

	return detectedRemote, detectedBranch, nil
}

// getUpstreamFromConfig gets upstream remote and branch from git config
func getUpstreamFromConfig(repo *git.Repository, currentBranch string) (remote, branch string, ok bool) {
	branchConfig, err := repo.Branch(currentBranch)
	if err != nil || branchConfig.Remote == "" {
		return "", "", false
	}

	remote = branchConfig.Remote
	branch = branchConfig.Merge.Short()

	if branch == "" {
		return "", "", false
	}

	return remote, branch, true
}

// autoDetectRemote auto-detects the remote
func autoDetectRemote(repo *git.Repository) (string, error) {
	remotes, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("failed to list remotes: %w", err)
	}

	var foundUpstream bool
	for _, remote := range remotes {
		name := remote.Config().Name
		if name == "origin" {
			return "origin", nil
		}
		if name == "upstream" {
			foundUpstream = true
		}
	}

	if foundUpstream {
		return "upstream", nil
	}

	return "", fmt.Errorf("no suitable remote found (origin or upstream)")
}
