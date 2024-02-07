// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
)

var logger = hclog.New(hclog.DefaultOptions)

var goModulesToUpdate = []goModuleType{
	baseLayerGoModule,
	resourceManagerGoModule,
}

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
	f.StringVar(&input.goSdkRepoPath, "go-sdk-repo-path", "", "--go-sdk-repo-path=../../../../go-azure-sdk")
	f.StringVar(&input.newSdkVersion, "new-sdk-version", "", "--new-sdk-version=1.4.0")
	f.StringVar(&input.outputFileName, "output-file", "", "--output-file=pr-description.txt")
	f.Parse(os.Args[1:])

	if err := input.validate(); err != nil {
		log.Fatalf("validating input: %+v", err)
	}

	absolutePathToAzureRMProvider, err := filepath.Abs(input.azurermRepoPath)
	if err != nil {
		log.Fatalf("determining absolute path to the AzureRM Provider at %q: %+v", input.azurermRepoPath, err)
	}
	input.azurermRepoPath = absolutePathToAzureRMProvider

	if input.goSdkRepoPath != "" {
		absolutePathToGoSdk, err := filepath.Abs(input.goSdkRepoPath)
		if err != nil {
			log.Fatalf("determining absolute path to the Go SDK at %q: %+v", input.goSdkRepoPath, err)
		}
		input.goSdkRepoPath = absolutePathToGoSdk
	}

	if err := run(context.Background(), input); err != nil {
		log.Fatalf(err.Error())
	}
}

type config struct {
	azurermRepoPath string
	goSdkRepoPath   string
	newSdkVersion   string
	outputFileName  string
}

func (c config) validate() error {
	if c.azurermRepoPath == "" {
		return fmt.Errorf("`--azurerm-repo-path` must be specified")
	}
	if c.newSdkVersion == "" {
		return fmt.Errorf("`--new-sdk-version` must be specified")
	}

	return nil
}

func run(ctx context.Context, input config) error {
	// Clone the repository
	// Determine what's changed in this
	//// git diff --name-only --diff-filter=A  v0.20220711.1181406...v0.20220712.1062733
	// ^ gives a list of paths which wants translating into `service[oldapi-newapi]`

	logger.Info(fmt.Sprintf("New SDK Version is %q", input.newSdkVersion))
	logger.Info(fmt.Sprintf("The `hashicorp/terraform-provider-azurerm` repository is located at %q", input.azurermRepoPath))

	if input.goSdkRepoPath != "" {
		logger.Info(fmt.Sprintf("The `hashicorp/go-azure-sdk` repository is located at %q", input.goSdkRepoPath))
	} else {
		logger.Info("A path to the `hashicorp/go-azure-sdk` repository was not provided - will clone on-demand")
	}

	if input.outputFileName != "" {
		logger.Info(fmt.Sprintf("Output File Name is %q", input.outputFileName))
	} else {
		logger.Info("No output file was specified so the PR description will only be output to the console")
	}

	// 1. Determine the current version of `hashicorp/go-azure-sdk` vendored into the Provider
	logger.Info("Determining the current version of `hashicorp/go-azure-sdk` being used..")
	oldSdkVersion, err := determineCurrentVersionOfGoAzureSDK(input.azurermRepoPath)
	if err != nil {
		return fmt.Errorf("determining the current version of `hashicorp/go-azure-sdk` being used in %q: %+v", input.azurermRepoPath, err)
	}
	logger.Info(fmt.Sprintf("Old SDK Version is %q", *oldSdkVersion))

	// 2. First determine the changes present in this version of the Go SDK
	// if there's no changes to the `resource-manager` or `sdk` folders, we can ignore it for now.
	logger.Info(fmt.Sprintf("Checking the changes between %q and %q of `hashicorp/go-azure-sdk`..", *oldSdkVersion, input.newSdkVersion))
	changes, err := determineChangesBetweenVersionsOfGoAzureSDK(ctx, *oldSdkVersion, input.newSdkVersion, input.goSdkRepoPath)
	if err != nil {
		return fmt.Errorf("determining the changes between version %q and %q of `hashicorp/go-azure-sdk`: %+v", *oldSdkVersion, input.newSdkVersion, err)
	}
	if !changes.hasChangesToSdk && !changes.hasChangesToResourceManager {
		logger.Info("No changes to either the SDK or Resource Manager - skipping updating")
		return nil
	}

	// 3. Update the version of `hashicorp/go-azure-sdk` used in `terraform-provider-azurerm`
	for _, moduleType := range goModulesToUpdate {
		logger.Info(fmt.Sprintf("Updating the %q Go module within `hashicorp/go-azure-sdk`..", string(moduleType)))
		if err := updateVersionOfGoAzureSDK(ctx, input.azurermRepoPath, moduleType, input.newSdkVersion); err != nil {
			return fmt.Errorf("updating the version of the %q Go Module within `hashicorp/go-azure-sdk`: %+v", string(moduleType), err)
		}
	}

	// 4. Then for each new Service/API Version:
	//   a. Try updating to the new API Version
	//   b. `go mod tidy && go mod vendor`
	//   c. `go test -v ./internal/services/{serviceName}/...
	//   d. Commit if it works - otherwise reset it and track it on a list
	logger.Debug("Attempting to update any existing Services present within the Provider..")
	// from the services we've got, work through and determine which API Versions exist for this service
	results := make([]updatedServiceSummary, 0)
	if changes.hasChangesToResourceManager {
		// sort the service names for consistency in output
		serviceNames := make([]string, 0)
		for serviceName := range changes.newServicesToApiVersions {
			serviceNames = append(serviceNames, serviceName)
		}
		sort.Strings(serviceNames)

		for _, serviceName := range serviceNames {
			availableApiVersions := changes.newServicesToApiVersions[serviceName]
			hasPendingChanges, err := directoryHasPendingChanges(input.azurermRepoPath)
			if err != nil {
				return fmt.Errorf("checking for pending changes in %q: %+v", input.azurermRepoPath, err)
			}
			if *hasPendingChanges {
				return fmt.Errorf("internal-error: working directory was not clean before starting service %q", serviceName)
			}

			logger.Info(fmt.Sprintf("Processing Service %q..", serviceName))
			apiVersionsCurrentlyUsedForService, err := determineApiVersionsCurrentlyUsedForService(input.azurermRepoPath, serviceName)
			if err != nil {
				return fmt.Errorf("determining the api versions currently used for service %q: %+v", serviceName, err)
			}
			for _, existingVersion := range *apiVersionsCurrentlyUsedForService {
				for _, newApiVersion := range availableApiVersions {
					if !shouldUpdateFrom(existingVersion, newApiVersion) {
						logger.Info(fmt.Sprintf("Skipping update from API Version %q to %q for Service %q", existingVersion, newApiVersion, serviceName))
						continue
					}

					logger.Info(fmt.Sprintf("Attempting to update API Version %q to %q for Service %q..", existingVersion, newApiVersion, serviceName))
					internalDirectory := path.Join(input.azurermRepoPath, "internal")
					if err := updateImportsWithinDirectory(serviceName, existingVersion, newApiVersion, internalDirectory); err != nil {
						return fmt.Errorf("updating the imports within %q: %+v", internalDirectory, err)
					}
					logger.Info(fmt.Sprintf("Updated the Imports for Service %q to use API Version %q rather than %q..", serviceName, newApiVersion, existingVersion))

					logger.Info("Running `go mod tidy` / `go mod vendor`..")
					goModTidyAndVendor(ctx, input.azurermRepoPath)

					logger.Debug("Checking for pending changes..")
					hasPendingChanges, err := directoryHasPendingChanges(input.azurermRepoPath)
					if err != nil {
						return fmt.Errorf("checking for pending changes in %q: %+v", input.azurermRepoPath, err)
					}
					if !*hasPendingChanges {
						logger.Info(fmt.Sprintf("No pending changes after attempting to update API Version %q to %q for Service %q - skipping", existingVersion, newApiVersion, serviceName))
						continue
					}

					logger.Info(fmt.Sprintf("Running `make test` within %q..", input.azurermRepoPath))
					if err := runMakeTest(ctx, input.azurermRepoPath); err != nil {
						results = append(results, updatedServiceSummary{
							serviceName:     serviceName,
							olderApiVersion: existingVersion,
							newApiVersion:   newApiVersion,
							error:           pointer.To(err.Error()),
						})

						// when the tests fail, we need to reset the working directory to ensure that there aren't any unstaged changes
						logger.Info("Resetting the working directory since `make test` failed..")
						if err := resetWorkingDirectory(ctx, input.azurermRepoPath); err != nil {
							return fmt.Errorf("resetting the working directory: %+v", err)
						}
						goModTidyAndVendor(ctx, input.azurermRepoPath)
						continue
					}

					logger.Debug("Committing changes..")
					message := fmt.Sprintf("dependencies: updating `%s` to API Version %q from %q", serviceName, newApiVersion, existingVersion)
					if err := stageAndCommitChanges(input.azurermRepoPath, message); err != nil {
						return fmt.Errorf("staging/committing changes when updating to API Version %q for %q: %+v", newApiVersion, serviceName, err)
					}

					logger.Info(fmt.Sprintf("Updated Service %q from %q to %q", serviceName, existingVersion, newApiVersion))
					results = append(results, updatedServiceSummary{
						serviceName:     serviceName,
						olderApiVersion: existingVersion,
						newApiVersion:   newApiVersion,
						error:           nil,
					})
				}
			}

			logger.Info(fmt.Sprintf("Processed Service %q.", serviceName))
		}
	}

	// 5. Build up a summary which can be used as a PR description
	logger.Debug("Building and outputting the PR description..")
	description := buildPullRequestDescription(results, input.newSdkVersion)
	if input.outputFileName != "" {
		logger.Info(fmt.Sprintf("Writing PR description to %q..", input.outputFileName))
		if err := os.WriteFile(input.outputFileName, []byte(description), 0644); err != nil {
			return fmt.Errorf("writing description to `%s`: %+v", input.outputFileName, err)
		}

		logger.Info(fmt.Sprintf("Processing completed - details written to %q", input.outputFileName))
	} else {
		logger.Info("Writing PR description to stdout since an output file was not specified")
		logger.Info("Summary of changes:")
		logger.Info(description)
	}

	return nil
}

func determineCurrentVersionOfGoAzureSDK(workingDirectory string) (*string, error) {
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

		// Since each of the Nested Go Modules within `hashicorp/go-azure-sdk` are versioned the same with a different
		// prefix, checking the base layer Go Module should be sufficient.
		if strings.EqualFold(v.Mod.Path, "github.com/hashicorp/go-azure-sdk/sdk") {
			return pointer.To(v.Mod.Version), nil
		}
	}

	return nil, fmt.Errorf("couldn't find the go module version for `github.com/hashicorp/go-azure-sdk/sdk` in %q", filePath)
}

func resetWorkingDirectory(ctx context.Context, workingDirectory string) error {
	repo, err := git.PlainOpen(workingDirectory)
	if err != nil {
		return fmt.Errorf("opening repository at %q: %+v", workingDirectory, err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("obtaining worktree: %+v", err)
	}

	logger.Trace(fmt.Sprintf("Performing `git reset --hard` in %s..", workingDirectory))
	opts := &git.ResetOptions{
		Mode: git.HardReset,
	}
	if err := worktree.Reset(opts); err != nil {
		return fmt.Errorf("performing reset: %+v", err)
	}

	// Whilst the git library exposes a `clean`, it doesn't support `clean -xdf`
	// so we'll have to shell out for that
	logger.Trace(fmt.Sprintf("Performing `git clean` in %s..", workingDirectory))
	args := []string{
		"clean",
		"-xdf",
	}
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = workingDirectory
	_ = cmd.Start()
	_ = cmd.Wait()
	return nil
}

func shouldUpdateFrom(existing string, new string) bool {
	if strings.EqualFold(existing, new) {
		logger.Trace("Versions are the same")
		return false
	}
	lowered := strings.ToLower(new)
	if strings.Contains(lowered, "alpha") {
		logger.Trace("Is an alpha..")
		return false
	}
	if strings.Contains(lowered, "beta") {
		logger.Trace("Is a beta..")
		return false
	}
	if strings.Contains(lowered, "preview") {
		logger.Trace("Is a preview..")
		return false
	}

	return new > existing
}

func runMakeTest(ctx context.Context, workingDirectory string) error {
	absPath, err := filepath.Abs(workingDirectory)
	if err != nil {
		return fmt.Errorf("determining absolute path for %q: %+v", workingDirectory, err)
	}

	logger.Trace("Running `make test`..")
	args := []string{
		"test",
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "make", args...)
	cmd.Dir = absPath
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	_ = cmd.Start()
	_ = cmd.Wait()
	if stderr.Len() > 0 {
		output := fmt.Sprintf(`stdout:
---
%s
---

stderr:
---
%s
---
`, stdout.String(), stderr.String())
		return fmt.Errorf("running `make test` in %q: %s", workingDirectory, output)
	}

	return nil
}

type updatedServiceSummary struct {
	serviceName     string
	olderApiVersion string
	newApiVersion   string
	error           *string
}

func (s updatedServiceSummary) successful() bool {
	return s.error == nil
}

func determineApiVersionsCurrentlyUsedForService(workingDirectory string, serviceName string) (*[]string, error) {
	absPath, err := filepath.Abs(workingDirectory)
	if err != nil {
		return nil, fmt.Errorf("determining absolute path for %q: %+v", workingDirectory, err)
	}

	logger.Debug(fmt.Sprintf("Determining the API Versions used for `hashicorp/go-azure-sdk` Service %q..", serviceName))
	servicesDirectory := path.Join(absPath, "internal", "services")
	imports, err := findImportsWithinDirectory(servicesDirectory)
	if err != nil {
		return nil, fmt.Errorf("finding imports within %q: %+v", servicesDirectory, err)
	}

	// now that we have a canonical list, unique these
	apiVersions := make(map[string]struct{}, 0)
	for _, line := range *imports {
		logger.Trace(fmt.Sprintf("Parsing import line %q..", line))

		// if there's an import alias, we need to remove that to determine the correct API version
		if !strings.HasPrefix(line, `"`) {
			// if the import is (force) imported (even if unused)
			if strings.HasPrefix(line, "_") {
				line = strings.TrimPrefix(line, "_")
				line = strings.TrimSpace(line)
			}

			if !strings.Contains(line, " ") {
				return nil, fmt.Errorf("the import line %q looks like an import alias but didn't parse in the format `alias \"importpath\"`", line)
			}
			components := strings.Split(line, " ")
			if len(components) != 2 {
				return nil, fmt.Errorf("expected the import alias to be in the format `alias \"path\"` but got %q which was %d segments: %+v", line, len(components), components)
			}
			line = strings.TrimSpace(components[1])
		}
		serviceImportPath := fmt.Sprintf("github.com/hashicorp/go-azure-sdk/resource-manager/%s/", serviceName)
		if !strings.Contains(line, serviceImportPath) {
			logger.Trace(fmt.Sprintf("Skipping line %q since it's not for this SDK..", line))
			continue
		}

		// pull out the api version, which is predictable
		line = strings.TrimPrefix(line, `"`)
		line = strings.TrimPrefix(line, serviceImportPath)
		line = strings.TrimSuffix(line, `"`)
		components := strings.Split(line, "/")
		apiVersion := components[0]
		logger.Trace(fmt.Sprintf("Found API Version %q from %q", apiVersion, line))
		apiVersions[apiVersion] = struct{}{}
	}

	out := make([]string, 0)
	for k := range apiVersions {
		out = append(out, k)
	}
	sort.Strings(out)
	return &out, nil
}

type changes struct {
	// hasChangesToResourceManager specifies that there are changes to the `resource-manager` directory
	hasChangesToResourceManager bool

	// hasChangesToSdk specifies that there are changes to the base layer/SDK
	hasChangesToSdk bool

	// newServicesToApiVersions contains a list of new API versions (value) for the given service (key)
	newServicesToApiVersions map[string][]string
}

func determineChangesBetweenVersionsOfGoAzureSDK(ctx context.Context, oldSDKVersion, newSDKVersion, goSdkRepositoryPath string) (*changes, error) {
	workingDirectory := goSdkRepositoryPath
	if goSdkRepositoryPath == "" {
		tempDirectory := os.TempDir()
		logger.Debug(fmt.Sprintf("Creating Temp Directory at %q..", tempDirectory))
		_ = os.MkdirAll(tempDirectory, 755)
		workingDirectory = path.Join(tempDirectory, "go-azure-sdk")
		defer func() {
			logger.Debug(fmt.Sprintf("Cleaning up the temp directory at %q..", workingDirectory))
			_ = os.RemoveAll(workingDirectory)
		}()

		logger.Debug(fmt.Sprintf("Cloning `hashicorp/go-azure-sdk` into %q..", workingDirectory))
		args := []string{
			"clone",
			"https://github.com/hashicorp/go-azure-sdk.git",
			workingDirectory,
		}
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = tempDirectory
		_ = cmd.Start()
		_ = cmd.Wait()
	}

	lines := make([]string, 0)
	for _, moduleType := range goModulesToUpdate {
		// Obtain the file paths which have changed
		// For now this should be sufficient to pull any changes and iterate over those e.g.:
		// > git diff --name-only v0.20220711.1181406...v0.20220712.1062733
		//
		// There's probably a performance enhancement here using the additional flags on
		// `--diff-filter` which supports Added (A) - but there's also Changed (C), Deleted (D)
		// and Renamed (R) - and can be used like so:
		// > git diff --name-only --diff-filter=A  v0.20220711.1181406...v0.20220712.1062733
		// however doing so will come with edge-cases, so it's simplest to pull all changes for now.
		logger.Debug(fmt.Sprintf("Determining the changes between %q and %q for %q..", oldSDKVersion, newSDKVersion, string(moduleType)))
		diffArgs := []string{
			"diff",
			"--name-only",
			fmt.Sprintf("%[1]s/%[2]s...%[1]s/%[3]s", string(moduleType), oldSDKVersion, newSDKVersion),
		}
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		diffCmd := exec.CommandContext(ctx, "git", diffArgs...)
		diffCmd.Dir = workingDirectory
		diffCmd.Stderr = &stderr
		diffCmd.Stdout = &stdout
		_ = diffCmd.Start()
		_ = diffCmd.Wait()
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("determining the changes between `hashicorp/go-azure-sdk` version %q and %q: %s", oldSDKVersion, newSDKVersion, stderr.String())
		}

		logger.Debug("Parsing changes from the Git Diff..")
		linesWithChanges := strings.Split(stdout.String(), "\n")
		lines = append(lines, linesWithChanges...)
	}

	parsed := parseChangesFromGitDiff(lines)
	return &parsed, nil
}

func parseChangesFromGitDiff(lines []string) changes {
	out := changes{
		hasChangesToResourceManager: false,
		hasChangesToSdk:             false,
		newServicesToApiVersions:    make(map[string][]string),
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "sdk/") {
			out.hasChangesToSdk = true
			continue
		}

		if strings.HasPrefix(line, "resource-manager/") {
			out.hasChangesToResourceManager = true

			// resource-manager/connectedvmware/2023-10-01/datastores/method_update.go
			split := strings.Split(line, "/")
			if len(split) < 3 {
				logger.Debug(fmt.Sprintf("Ignoring line %q - not in the format `resource-manager/{serviceName}/{apiVersion}", line))
				continue
			}

			serviceName := split[1]
			apiVersion := split[2]
			existing, ok := out.newServicesToApiVersions[serviceName]
			if !ok {
				existing = make([]string, 0)
			}
			containsExistingValue := false
			for _, val := range existing {
				if val == apiVersion {
					containsExistingValue = true
					break
				}
			}
			if !containsExistingValue {
				existing = append(existing, apiVersion)
			}

			out.newServicesToApiVersions[serviceName] = existing

			continue
		}
	}
	return out
}

type goModuleType string

const (
	// baseLayerGoModule represents the Go Module containing the base layer in `hashicorp/go-azure-sdk`
	// this is the Go Module `github.com/hashicorp/go-azure-sdk/sdk`
	baseLayerGoModule goModuleType = "sdk"

	// resourceManagerGoModule represents the Go Module containing the Resource Manager SDK in
	// `hashicorp/go-azure-sdk` - which is the Go Module `github.com/hashicorp/go-azure-sdk/resource-manager`
	resourceManagerGoModule goModuleType = "resource-manager"
)

func updateVersionOfGoAzureSDK(ctx context.Context, workingDirectory string, moduleType goModuleType, newApiVersion string) error {
	logger.Debug(fmt.Sprintf("Updating the version of `go-azure-sdk`'s %q Go Module to %q", string(moduleType), newApiVersion))
	args := []string{
		"get",
		fmt.Sprintf("github.com/hashicorp/go-azure-sdk/%s@%s", string(moduleType), newApiVersion),
	}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = workingDirectory
	_ = cmd.Start()
	_ = cmd.Wait()
	logger.Debug(fmt.Sprintf("Updated the version of `go-azure-sdk` to %q", newApiVersion))

	logger.Debug("Vendoring the changes..")
	goModTidyAndVendor(ctx, workingDirectory)

	logger.Debug("Committing the changes..")
	message := fmt.Sprintf("dependencies: updating to version `%s` of `github.com/hashicorp/go-azure-sdk/%s`", newApiVersion, string(moduleType))
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

func directoryHasPendingChanges(workingDirectory string) (*bool, error) {
	repo, err := git.PlainOpen(workingDirectory)
	if err != nil {
		return nil, fmt.Errorf("opening %q: %+v", workingDirectory, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("opening worktree at %q: %+v", workingDirectory, err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("obtaining status for %q: %+v", workingDirectory, err)
	}

	return pointer.To(!status.IsClean()), nil
}

func updateImportsWithinDirectory(serviceName string, oldApiVersion string, newApiVersion string, workingDirectory string) error {
	absPath, err := filepath.Abs(workingDirectory)
	if err != nil {
		return fmt.Errorf("obtaining absolute path for %q: %+v", workingDirectory, err)
	}

	// because we want to process all nested directories, we need to first pull out a complete list of directories
	fileSet := token.NewFileSet()
	nestedDirectories := findDirectoriesNestedWithin(absPath)

	// over which we can then iterate to get a list of files within that directory
	for _, directory := range nestedDirectories {
		logger.Trace(fmt.Sprintf("Processing directory %q..", directory))
		files, err := parser.ParseDir(fileSet, directory, func(info fs.FileInfo) bool {
			return true
		}, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parsing files within %q: %+v", directory, err)
		}
		// to be able to update the imports for said directory
		for pkgName, pkg := range files {
			logger.Trace(fmt.Sprintf("Processing Go Package %q", pkgName))
			for fileName, file := range pkg.Files {
				logger.Trace(fmt.Sprintf("Updating imports for File %q..", fileName))
				updateImportsForFile(fileSet, file, serviceName, oldApiVersion, newApiVersion)

				var buf bytes.Buffer
				if err = format.Node(&buf, fileSet, file); err != nil {
					return fmt.Errorf("error formatting new code: %w", err)
				}
				_ = os.WriteFile(fileName, buf.Bytes(), 0644)
			}
		}
		logger.Trace(fmt.Sprintf("Processed directory %q.", directory))
	}
	return nil
}

func updateImportsForFile(fileSet *token.FileSet, file *ast.File, serviceName string, oldApiVersion string, newApiVersion string) {
	importLineForPreviousApiVersion := fmt.Sprintf("github.com/hashicorp/go-azure-sdk/resource-manager/%s/%s", serviceName, oldApiVersion)
	importLineForNewApiVersion := fmt.Sprintf("github.com/hashicorp/go-azure-sdk/resource-manager/%s/%s", serviceName, newApiVersion)

	// first update the imports themselves
	existingImports := astutil.Imports(fileSet, file)
	aliasesToReplace := make(map[string]string)
	for _, val := range existingImports {
		for _, item := range val {
			logger.Trace(fmt.Sprintf("Processing Import %q", item.Path.Value))
			existingImportLine := item.Path.Value
			if !strings.Contains(existingImportLine, importLineForPreviousApiVersion) {
				continue
			}

			updatedImportLine := strings.Replace(existingImportLine, importLineForPreviousApiVersion, importLineForNewApiVersion, 1)
			logger.Trace(fmt.Sprintf("Updating Import URI from %q to %q", existingImportLine, updatedImportLine))
			item.Path.Value = updatedImportLine

			// if we're importing the meta client (e.g. the api version directly) then we also need to update the alias
			importsMetaClient := strings.ReplaceAll(existingImportLine, "\"", "") == importLineForPreviousApiVersion
			if importsMetaClient && item.Name != nil {
				if existingAlias := item.Name.Name; existingAlias != "" {
					updatedAlias := strings.ToLower(fmt.Sprintf("%s_%s", serviceName, strings.ReplaceAll(newApiVersion, "-", "_")))

					logger.Trace(fmt.Sprintf("Updating Import Alias from %q to %q", existingAlias, updatedAlias))
					aliasesToReplace[existingAlias] = updatedAlias
					item.Name.Name = updatedAlias
				}
			}

			// finally, remove any comments which will be stragglers/lintignores which shouldn't be present
			if item.Comment != nil {
				item.Comment.List = []*ast.Comment{}
			}
		}
	}

	// then update any references to the aliases we've updated
	ast.Inspect(file, func(n ast.Node) bool {
		v, ok := n.(*ast.Ident)
		if ok {
			for alias, replacement := range aliasesToReplace {
				if v.Name == alias {
					v.Name = replacement
				}
			}
		}

		return true
	})
}

func findImportsWithinDirectory(workingDirectory string) (*[]string, error) {
	absPath, err := filepath.Abs(workingDirectory)
	if err != nil {
		return nil, fmt.Errorf("obtaining absolute path for %q: %+v", workingDirectory, err)
	}

	imports := make(map[string]struct{})
	nestedDirectories := findDirectoriesNestedWithin(absPath)
	for _, directory := range nestedDirectories {
		logger.Trace(fmt.Sprintf("Processing directory %q..", absPath))

		fileSet := token.NewFileSet()
		files, err := parser.ParseDir(fileSet, directory, func(info fs.FileInfo) bool {
			return true
		}, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parsing files within %q: %+v", directory, err)
		}

		for pkgName, pkg := range files {
			logger.Trace(fmt.Sprintf("Processing Go Package %q", pkgName))
			for fileName, file := range pkg.Files {
				logger.Trace(fmt.Sprintf("Finding imports within File %q..", fileName))
				importsInFile := findImportsWithinFile(fileSet, file)
				for _, item := range importsInFile {
					imports[item] = struct{}{}
				}
			}
		}
		logger.Trace(fmt.Sprintf("Processed directory %q.", directory))
	}

	out := make([]string, 0)
	for k := range imports {
		out = append(out, k)
	}
	sort.Strings(out)
	return &out, nil
}

func findDirectoriesNestedWithin(workingDirectory string) []string {
	// because we want to process all nested directories, we need to first pull out a complete list of directories
	directories := make([]string, 0)
	_ = filepath.WalkDir(workingDirectory, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})
	sort.Strings(directories)
	return directories
}

func findImportsWithinFile(fileSet *token.FileSet, file *ast.File) []string {
	existingImports := astutil.Imports(fileSet, file)
	imports := make(map[string]struct{})
	for _, val := range existingImports {
		for _, item := range val {
			logger.Trace(fmt.Sprintf("Processing Import %q", item.Path.Value))
			existingImportLine := item.Path.Value
			if !strings.Contains(existingImportLine, "github.com/hashicorp/go-azure-sdk/resource-manager/") {
				continue
			}

			imports[item.Path.Value] = struct{}{}
		}
	}

	out := make([]string, 0)
	for k := range imports {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func buildPullRequestDescription(input []updatedServiceSummary, newSdkVersion string) string {
	succeeded := make([]updatedServiceSummary, 0)
	failed := make([]updatedServiceSummary, 0)
	for _, v := range input {
		if v.successful() {
			succeeded = append(succeeded, v)
			continue
		}

		failed = append(failed, v)
		continue
	}

	succeededLines := make([]string, 0)
	if len(succeeded) > 0 {
		// we're intentionally omitting a header tag here to make it clearer when something has failed
		succeededLines = append(succeededLines, fmt.Sprintf(`This updates the following Services and API Versions:`))
		for _, result := range succeeded {
			succeededLines = append(succeededLines, fmt.Sprintf("* The service `%s` was updated to API Version `%s` (from `%s`).", result.serviceName, result.newApiVersion, result.olderApiVersion))
		}
	}

	failedLines := make([]string, 0)
	if len(failed) > 0 {
		failedLines = append(failedLines, `## FAILED - API Versions`)
		failedLines = append(failedLines, fmt.Sprintf(`The following new API Versions are available but had compile-time errors when updating:`))
		for i, result := range failed {
			failedLines = append(failedLines, fmt.Sprintf("* The service `%s` - updating to API Version `%s` from `%s`.", result.serviceName, result.newApiVersion, result.olderApiVersion))
			failedLines = append(failedLines, fmt.Sprintf("```\n%s\n```", *result.error))
			if i != len(failed)-1 {
				failedLines = append(failedLines, "\n---")
			}
		}
	}

	otherLines := make([]string, 0)
	if len(succeeded) == 0 && len(failed) == 0 {
		otherLines = append(otherLines, "This version contains no new API Versions - so no upgrades were attempted.")
	}

	lines := []string{
		fmt.Sprintf("This PR updates the version of `hashicorp/go-azure-sdk` to `%s`.", newSdkVersion),
		strings.Join(succeededLines, "\n"),
		strings.Join(failedLines, "\n"),
		strings.Join(otherLines, "\n"),
	}
	return strings.Join(lines, "\n\n")
}
