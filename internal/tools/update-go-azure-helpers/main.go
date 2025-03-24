// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/mod/modfile"
)

var logger = hclog.New(hclog.DefaultOptions)

func main() {
	// NOTE: this tool has a number of assumptions both about the environment this is being run in, and
	// about the design of this repository - namely that an `internal` folder exists and that Git is configured.

	if v := os.Getenv("LOG_LEVEL"); v != "" {
		level := hclog.LevelFromString(v)
		logger.SetLevel(level)
	}

	input := config{}

	f := flag.NewFlagSet("update-go-azure-sdk", flag.PanicOnError)
	f.StringVar(&input.azurermRepoPath, "azurerm-repo-path", "", "--azurerm-repo-path=../../../")
	f.StringVar(&input.azureHelpersRepoPath, "azure-helpers-repo-path", "", "--azure-helpers-repo-path=../../../../go-azure-helpers")
	f.StringVar(&input.newHelpersVersion, "new-helpers-version", "", "--new-helpers-version=1.4.0")
	f.Parse(os.Args[1:])

	if err := input.validate(); err != nil {
		log.Fatalf("validating input: %+v", err)
	}

	absolutePathToAzureRMProvider, err := filepath.Abs(input.azurermRepoPath)
	if err != nil {
		log.Fatalf("determining absolute path to the AzureRM Provider at %q: %+v", input.azurermRepoPath, err)
	}
	input.azurermRepoPath = absolutePathToAzureRMProvider

	if input.azureHelpersRepoPath != "" {
		absolutePathToHelpers, err := filepath.Abs(input.azureHelpersRepoPath)
		if err != nil {
			log.Fatalf("determining absolute path to the Go Azure Helpers at %q: %+v", input.azureHelpersRepoPath, err)
		}
		input.azureHelpersRepoPath = absolutePathToHelpers
	}

	if err := run(context.Background(), input); err != nil {
		log.Fatalf(err.Error())
	}
}

type config struct {
	azurermRepoPath      string
	azureHelpersRepoPath string
	newHelpersVersion    string
}

func (c config) validate() error {
	if c.azurermRepoPath == "" {
		return fmt.Errorf("`--azurerm-repo-path` must be specified")
	}
	if c.newHelpersVersion == "" {
		return fmt.Errorf("`--new-helpers-version` must be specified")
	}

	return nil
}

func run(ctx context.Context, input config) error {
	logger.Info(fmt.Sprintf("New Go Azure Helpers Version is %q", input.newHelpersVersion))
	logger.Info(fmt.Sprintf("The `hashicorp/terraform-provider-azurerm` repository is located at %q", input.azurermRepoPath))

	if input.azureHelpersRepoPath != "" {
		logger.Info(fmt.Sprintf("The `hashicorp/go-azure-helpers` repository is located at %q", input.azureHelpersRepoPath))
	} else {
		logger.Info("A path to the `hashicorp/go-azure-helpers` repository was not provided - will clone on-demand")
	}

	// 1. Determine the current version of `hashicorp/go-azure-helpers` vendored into the Provider
	logger.Info("Determining the current version of `hashicorp/go-azure-helpers` being used..")
	oldHelpersVersion, err := determineCurrentVersionOfGoAzureHelpers(input.azurermRepoPath)
	if err != nil {
		return fmt.Errorf("determining the current version of `hashicorp/go-azure-helpers` being used in %q: %+v", input.azurermRepoPath, err)
	}
	logger.Info(fmt.Sprintf("Old Go Azure Helpers Version is %q", *oldHelpersVersion))

	// 2. Update the version of `hashicorp/go-azure-helpers` used in `terraform-provider-azurerm`
	logger.Info(fmt.Sprintf("Updating `hashicorp/go-azure-helpers`.."))
	if err := updateVersionOfGoAzureHelpers(ctx, input.azurermRepoPath, input.newHelpersVersion); err != nil {
		return fmt.Errorf("updating the version of `hashicorp/go-azure-helpers`: %+v", err)
	}

	return nil
}

func determineCurrentVersionOfGoAzureHelpers(workingDirectory string) (*string, error) {
	filePath := path.Join(workingDirectory, "go.mod")
	logger.Trace(fmt.Sprintf("Parsing the go.mod at %q..", filePath))
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading go.mod at %q: %+v", filePath, err)
	}

	module, err := modfile.Parse(filePath, file, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing the go.mod at %q: %+v", filePath, err)
	}

	for _, v := range module.Require {
		if v == nil {
			continue
		}

		// Since each of the Nested Go Modules within `hashicorp/go-azure-helpers` are versioned the same with a different
		// prefix, checking the base layer Go Module should be sufficient.
		if strings.EqualFold(v.Mod.Path, "github.com/hashicorp/go-azure-helpers") {
			return pointer.To(v.Mod.Version), nil
		}
	}

	return nil, fmt.Errorf("couldn't find the go module version for `github.com/hashicorp/go-azure-helpers` in %q", filePath)
}

func updateVersionOfGoAzureHelpers(ctx context.Context, workingDirectory string, newApiVersion string) error {
	logger.Debug(fmt.Sprintf("Updating the version of `go-azure-helpers` to %q", newApiVersion))
	args := []string{
		"get",
		fmt.Sprintf("github.com/hashicorp/go-azure-helpers/%s", newApiVersion),
	}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = workingDirectory
	_ = cmd.Start()
	_ = cmd.Wait()
	logger.Debug(fmt.Sprintf("Updated the version of `go-azure-helpers` to %q", newApiVersion))

	logger.Debug("Vendoring the changes..")
	goModTidyAndVendor(ctx, workingDirectory)

	logger.Debug("Committing the changes..")
	message := fmt.Sprintf("dependencies: updating to version `%s` of `github.com/hashicorp/go-azure-helpers`", newApiVersion)
	if err := stageAndCommitChanges(workingDirectory, message); err != nil {
		return fmt.Errorf("staging/committing changes to %q: %+v", workingDirectory, err)
	}
	return nil
}

func goModTidyAndVendor(ctx context.Context, workingDirectory string) {
	logger.Debug(fmt.Sprintf("Running `go mod tidy` in %q..", workingDirectory))
	tidyArgs := []string{
		"mod",
		"tidy",
	}
	tidyCmd := exec.CommandContext(ctx, "go", tidyArgs...)
	tidyCmd.Dir = workingDirectory
	_ = tidyCmd.Start()
	_ = tidyCmd.Wait()
	logger.Debug(fmt.Sprintf("Run `go mod tidy` in %q.", workingDirectory))

	logger.Debug(fmt.Sprintf("Running `go mod vendor` in %q..", workingDirectory))
	vendorArgs := []string{
		"mod",
		"vendor",
	}
	vendorCmd := exec.CommandContext(ctx, "go", vendorArgs...)
	vendorCmd.Dir = workingDirectory
	_ = vendorCmd.Start()
	_ = vendorCmd.Wait()
	logger.Debug(fmt.Sprintf("Run `go mod vendor` in %q.", workingDirectory))
	return
}

func stageAndCommitChanges(workingDirectory string, message string) error {
	repo, err := git.PlainOpen(workingDirectory)
	if err != nil {
		return fmt.Errorf("opening %q: %+v", workingDirectory, err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("opening worktree: %+v", err)
	}

	// Whilst `Commit` has an add below it doesn't capture all changes, hence we need
	// to explicitly stage all changes (including the vendored changes)
	logger.Trace(fmt.Sprintf("Staging all changes in %q..", workingDirectory))
	addOpts := &git.AddOptions{
		// NOTE: Glob is needed to add all untracked files too, else the vendor directory isn't committed
		All:  true,
		Glob: "**",
	}
	if err := worktree.AddWithOptions(addOpts); err != nil {
		return fmt.Errorf("staging changes: %+v", err)
	}

	logger.Trace(fmt.Sprintf("Committing all changes in %q..", workingDirectory))
	opts := &git.CommitOptions{
		// locally the author/committer info comes from the `.gitconfig`
	}
	if os.Getenv("RUNNING_IN_AUTOMATION") != "" {
		// however in automation lets hardcode this
		opts.Author = &object.Signature{
			Name:  os.Getenv("GIT_COMMIT_USERNAME"),
			Email: "",
			When:  time.Now(),
		}
		opts.Committer = &object.Signature{
			Name:  os.Getenv("GIT_COMMIT_USERNAME"),
			Email: "",
			When:  time.Now(),
		}
	}

	hash, err := worktree.Commit(message, opts)
	if err != nil {
		return fmt.Errorf("commiting changes to %q: %+v", workingDirectory, err)
	}

	logger.Info(fmt.Sprintf("Committed as %q", hash))
	return nil
}
