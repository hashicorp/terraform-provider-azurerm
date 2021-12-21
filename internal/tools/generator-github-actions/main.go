package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	directoryContainingServicePackages    = "./internal/services/"
	golangVersion                         = "1.17.5"
	filePermissions                       = os.FileMode(0644)
	githubActionsPath                     = "./.github/workflows/"
	fileContainingNewServicesGithubAction = "./.github/workflows/tests-and-linting-for-new-services.yaml"
	fileContainingListOfServicePackages   = "./scripts/.service-packages"
)

var debug = false

func main() {
	updatingServicesListAfterMerge := os.Getenv("UPDATE_SERVICES_AFTER_MERGE") != ""
	debug = os.Getenv("DEBUG") != ""

	if updatingServicesListAfterMerge {
		if debug {
			log.Printf("[DEBUG] Updating the Services list..")
		}
		if err := updateListOfServicePackages(); err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		if debug {
			log.Printf("[DEBUG] Generating a Github Action for each Service Package..")
		}
		if err := generateGithubActionsForServicePackages(); err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func generateGithubActionsForServicePackages() error {
	workflowsPath, err := filepath.Abs(githubActionsPath)
	if err != nil {
		return fmt.Errorf("determining absolute path for %q: %+v", githubActionsPath, err)
	}

	servicesPath, err := filepath.Abs(directoryContainingServicePackages)
	if err != nil {
		return fmt.Errorf("determining absolute path for %q: %+v", githubActionsPath, err)
	}

	// first remove any of the previously generated Github actions
	if err := deleteGeneratedGithubActions(workflowsPath); err != nil {
		return fmt.Errorf("deleting generated github actions: %+v", err)
	}

	// then output one Github action per Service Package
	servicePackages, err := getServicePackages(servicesPath)
	if err != nil {
		return fmt.Errorf("listing service packages: %+v", err)
	}
	for _, service := range *servicePackages {
		fileName := fmt.Sprintf("generated-service-%s.yaml", strings.ToLower(service))
		filePath := path.Join(workflowsPath, fileName)

		if debug {
			log.Printf("[DEBUG] Writing %q", fileName)
		}
		config := githubActionForServicePackage(service)
		os.WriteFile(filePath, []byte(config), filePermissions)
	}

	return nil
}

func updateListOfServicePackages() error {
	servicesPath, err := filepath.Abs(directoryContainingServicePackages)
	if err != nil {
		return fmt.Errorf("determining absolute path for %q: %+v", githubActionsPath, err)
	}

	servicePackagesFilePath, err := filepath.Abs(fileContainingListOfServicePackages)
	if err != nil {
		return fmt.Errorf("getting absolute path for %q: %+v", fileContainingListOfServicePackages, err)
	}
	services, err := getServicePackages(servicesPath)
	if err != nil {
		return fmt.Errorf("getting service packages: %+v", err)
	}
	fileContents := strings.Join(*services, "\n")

	// first update the list of known service packages
	if debug {
		log.Printf("[DEBUG] Updating the file containing the known service packages..")
	}
	os.Remove(servicePackagesFilePath)
	os.WriteFile(servicePackagesFilePath, []byte(fileContents), filePermissions)

	// then update the Github Action to exclude those when checking for new Service Packages to be run
	log.Printf("[DEBUG] Updating the new services Github Action..")
	githubActionPath, err := filepath.Abs(fileContainingNewServicesGithubAction)
	if err != nil {
		return fmt.Errorf("getting absolute path for %q: %+v", githubActionPath, err)
	}
	fileContents = githubActionForNewServicesPackages(*services)
	if debug {
	}
	os.Remove(githubActionPath)
	os.WriteFile(githubActionPath, []byte(fileContents), filePermissions)

	return nil
}

func deleteGeneratedGithubActions(workflowsDirectory string) error {
	files, err := ioutil.ReadDir(workflowsDirectory)
	if err != nil {
		return fmt.Errorf("listing directory %q: %+v", workflowsDirectory, err)
	}

	for _, file := range files {
		filename := strings.ToLower(file.Name())
		if strings.HasPrefix(filename, "generated-") {
			if debug {
				log.Printf("[DEBUG] Removing existing file %q", file.Name())
			}
			filePath := path.Join(workflowsDirectory, file.Name())
			os.Remove(filePath)
		}
	}

	return nil
}

func getServicePackages(servicesDirectory string) (*[]string, error) {
	services := make([]string, 0)

	files, err := ioutil.ReadDir(servicesDirectory)
	if err != nil {
		return nil, fmt.Errorf("listing directory %q: %+v", servicesDirectory, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		services = append(services, file.Name())
	}

	sort.Strings(services)

	return &services, nil
}

func githubActionForServicePackage(name string) string {
	// NOTE: this intentionally omits `./github/workflows/**` from the list
	// this is to avoid running hundreds of Github Actions when we bump the
	// Go version (which presuming the provider builds, is sufficient)
	return fmt.Sprintf(`---
name: Unit Tests and Linting for Service Package %[1]q
on:
  pull_request:
    types: ['opened', 'synchronize']
    paths:
      - 'internal/services/%[1]s/**'

concurrency:
  group: 'unit-${{ github.head_ref }}'
  cancel-in-progress: true

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '%[2]s'
      - run: bash scripts/gogetcookie.sh
      - run: make tools
      - run: bash scripts/service-package.sh %[1]q
`, name, golangVersion)
}

func githubActionForNewServicesPackages(services []string) string {
	paths := make([]string, 0)
	for _, service := range services {
		paths = append(paths, fmt.Sprintf("      - '!internal/services/%[1]s/**'", service))
	}

	return fmt.Sprintf(`---
name: Unit Tests and Linting for New Service Packages
on:
  pull_request:
    types: ['opened', 'synchronize']
    paths:
      - 'internal/services/**'
%[1]s

concurrency:
  group: 'unit-${{ github.head_ref }}'
  cancel-in-progress: true

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '%[2]s'
      - run: bash scripts/gogetcookie.sh
      - run: make tools
      - run: bash scripts/run-new-service-packages.sh
`, strings.Join(paths, "\n"), golangVersion)
}
