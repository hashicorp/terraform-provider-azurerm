package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk" // nolint: typecheck
)

// Packages in this list are deprecated and cannot be run due to breaking API changes
// this should only be used as a last resort - since all acctests should run nightly.
var packagesToSkip = map[string]struct{}{
	"devspace": {},

	// force-deprecated and will be removed by Azure on 2021-04-28
	// new clusters cannot be cretaed - so there's no point trying to run the tests
	"servicefabricmesh": {},
}

func main() {
	filePath := flag.String("path", "", "The relative path to the root directory")
	showHelp := flag.Bool("help", false, "Display this message")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	generators := []generator{
		githubLabelsGenerator{},
		githubIssueLabelsGenerator{},
		teamCityServicesListGenerator{},
		websiteCategoriesGenerator{},
	}
	for _, value := range generators {
		outputFile := value.outputPath(*filePath)
		if err := value.run(outputFile, packagesToSkip); err != nil {
			panic(err)
		}
	}
}

type generator interface {
	outputPath(rootDirectory string) string
	run(outputFileName string, packagesToSkip map[string]struct{}) error
}

const githubLabelsTemplate = `# NOTE: this file is generated via 'make generate'
dependencies:
  - go.mod
  - go.sum
  - vendor/**/*
documentation:
  - website/**/*
tooling:
  - internal/tools/**/*
`

type githubLabelsGenerator struct{}

func (githubLabelsGenerator) outputPath(rootDirectory string) string {
	return fmt.Sprintf("%s/.github/labeler-pull-request-triage.yml", rootDirectory)
}

func (githubLabelsGenerator) run(outputFileName string, _ map[string]struct{}) error {
	packagesToLabel := make(map[string]string)
	// combine and unique these
	for _, service := range provider.SupportedTypedServices() {
		v, ok := service.(sdk.TypedServiceRegistrationWithAGitHubLabel)
		if !ok {
			// skipping since this doesn't implement the label interface
			continue
		}

		info := reflect.TypeOf(service)
		packageSegments := strings.Split(info.PkgPath(), "/")
		packageName := packageSegments[len(packageSegments)-1]
		packagesToLabel[packageName] = v.AssociatedGitHubLabel()
	}
	for _, service := range provider.SupportedUntypedServices() {
		v, ok := service.(sdk.UntypedServiceRegistrationWithAGitHubLabel)
		if !ok {
			// skipping since this doesn't implement the label interface
			continue
		}

		info := reflect.TypeOf(service)
		packageSegments := strings.Split(info.PkgPath(), "/")
		packageName := packageSegments[len(packageSegments)-1]
		packagesToLabel[packageName] = v.AssociatedGitHubLabel()
	}

	// labels can be present in more than one package, so we need to group them
	labelsToPackages := make(map[string][]string)
	for pkg, label := range packagesToLabel {
		existing, ok := labelsToPackages[label]
		if !ok {
			existing = []string{}
		}

		existing = append(existing, pkg)
		labelsToPackages[label] = existing
	}

	sortedLabels := make([]string, 0)
	for k := range labelsToPackages {
		sortedLabels = append(sortedLabels, k)
	}
	sort.Strings(sortedLabels)

	output := strings.TrimSpace(githubLabelsTemplate)
	for _, labelName := range sortedLabels {
		pkgs := labelsToPackages[labelName]

		// for consistent generation
		sort.Strings(pkgs)

		out := []string{
			fmt.Sprintf("%[1]s:", labelName),
		}
		for _, pkg := range pkgs {
			out = append(out, fmt.Sprintf("  - internal/services/%[1]s/**/*", pkg))
		}

		out = append(out, "")
		output += fmt.Sprintf("\n%s", strings.Join(out, "\n"))
	}

	return writeToFile(outputFileName, output)
}

type teamCityServicesListGenerator struct{}

func (teamCityServicesListGenerator) outputPath(rootDirectory string) string {
	return fmt.Sprintf("%s/.teamcity/components/generated/services.kt", rootDirectory)
}

func (teamCityServicesListGenerator) run(outputFileName string, packagesToSkip map[string]struct{}) error {
	template := `// NOTE: this is Generated from the Service Definitions - manual changes will be lost
//       to re-generate this file, run 'make generate' in the root of the repository
var services = mapOf(
%s
)`
	items := make([]string, 0)

	services := make(map[string]string)
	serviceNames := make([]string, 0)

	// combine and unique these
	for _, service := range provider.SupportedTypedServices() {
		info := reflect.TypeOf(service)
		packageSegments := strings.Split(info.PkgPath(), "/")
		packageName := packageSegments[len(packageSegments)-1]
		serviceName := service.Name()

		// Service Registrations are reused across Typed and Untyped Services now
		if _, exists := services[serviceName]; exists {
			continue
		}

		services[serviceName] = packageName
		serviceNames = append(serviceNames, serviceName)
	}
	for _, service := range provider.SupportedUntypedServices() {
		info := reflect.TypeOf(service)
		packageSegments := strings.Split(info.PkgPath(), "/")
		packageName := packageSegments[len(packageSegments)-1]
		serviceName := service.Name()

		// Service Registrations are reused across Typed and Untyped Services now
		if _, exists := services[serviceName]; exists {
			continue
		}

		services[serviceName] = packageName
		serviceNames = append(serviceNames, serviceName)
	}

	// then ensure these are sorted so they're alphabetical
	sort.Strings(serviceNames)
	for _, serviceName := range serviceNames {
		packageName := services[serviceName]
		if _, shouldSkip := packagesToSkip[packageName]; shouldSkip {
			continue
		}

		item := fmt.Sprintf("        %q to %q", packageName, serviceName)
		items = append(items, item)
	}

	formatted := fmt.Sprintf(template, strings.Join(items, ",\n"))
	return writeToFile(outputFileName, formatted)
}

type websiteCategoriesGenerator struct{}

func (websiteCategoriesGenerator) outputPath(rootDirectory string) string {
	return fmt.Sprintf("%s/website/allowed-subcategories", rootDirectory)
}

func (websiteCategoriesGenerator) run(outputFileName string, _ map[string]struct{}) error {
	websiteCategories := make([]string, 0)

	// get a distinct list
	for _, service := range provider.SupportedTypedServices() {
		for _, category := range service.WebsiteCategories() {
			if contains(websiteCategories, category) {
				continue
			}

			websiteCategories = append(websiteCategories, category)
		}
	}
	for _, service := range provider.SupportedUntypedServices() {
		for _, category := range service.WebsiteCategories() {
			if contains(websiteCategories, category) {
				continue
			}

			websiteCategories = append(websiteCategories, category)
		}
	}

	// sort them
	sort.Strings(websiteCategories)

	fileContents := strings.Join(websiteCategories, "\n")
	return writeToFile(outputFileName, fileContents)
}

const githubIssueLabelsTemplate = `# NOTE: this file is generated via 'make generate'
bug:
  - 'panic:'
crash:
  - 'panic:'
`

type githubIssueLabelsGenerator struct{}

func (g githubIssueLabelsGenerator) outputPath(rootDirectory string) string {
	return fmt.Sprintf("%s/.github/labeler-issue-triage.yml", rootDirectory)
}

func (githubIssueLabelsGenerator) run(outputFileName string, _ map[string]struct{}) error {
	labelToNames := make(map[string][]string)

	for _, service := range provider.SupportedTypedServices() {

		v, ok := service.(sdk.TypedServiceRegistrationWithAGitHubLabel)
		if !ok {
			// skipping since this doesn't implement the label interface
			continue
		}

		var names []string
		for _, resource := range service.Resources() {
			names = append(names, resource.ResourceType())

		}

		for _, ds := range service.DataSources() {
			names = append(names, ds.ResourceType())
		}

		appendToSliceWithinMap(labelToNames, names, v.AssociatedGitHubLabel())

	}
	for _, service := range provider.SupportedUntypedServices() {
		v, ok := service.(sdk.UntypedServiceRegistrationWithAGitHubLabel)
		service.SupportedResources()
		if !ok {
			// skipping since this doesn't implement the label interface
			continue
		}

		var names []string
		for resourceName := range service.SupportedResources() {
			if resourceName != "" {
				names = append(names, resourceName)
			}
		}

		for dsName := range service.SupportedDataSources() {
			if dsName != "" {
				names = append(names, dsName)
			}
		}

		appendToSliceWithinMap(labelToNames, names, v.AssociatedGitHubLabel())

	}

	sortedLabels := make([]string, 0)
	for k := range labelToNames {
		sortedLabels = append(sortedLabels, k)
	}
	sort.Strings(sortedLabels)

	output := strings.TrimSpace(githubIssueLabelsTemplate)
	for _, labelName := range sortedLabels {

		prefix := longestCommonPrefix(labelToNames[labelName])
		if prefix == "azurerm_" {
			continue
		}

		out := []string{
			fmt.Sprintf("%[1]s:", labelName),
		}

		out = append(out, fmt.Sprintf("  - '### (|New or )Affected Resource\\(s\\)\\/Data Source\\(s\\)((.|\\n)*)%s((.|\\n)*)###'", prefix))

		out = append(out, "")
		output += fmt.Sprintf("\n%s", strings.Join(out, "\n"))
	}

	return writeToFile(outputFileName, output)
}

func writeToFile(filePath string, contents string) error {
	outputPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// output that string to the file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	if os.IsExist(err) {
		os.Remove(outputPath)
		file, err = os.Create(outputPath)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	_, _ = file.WriteString(contents)

	return file.Sync()
}

func contains(input []string, value string) bool {
	for _, v := range input {
		if v == value {
			return true
		}
	}

	return false
}

func appendToSliceWithinMap(sliceMap map[string][]string, slice []string, key string) map[string][]string {
	if _, ok := sliceMap[key]; ok {
		sliceMap[key] = append(slice, sliceMap[key]...)
	} else {
		sliceMap[key] = slice
	}
	return sliceMap
}

func longestCommonPrefix(strings []string) string {
	longestPrefix := ""
	end := false

	if len(strings) > 0 {

		sort.Strings(strings)
		first := strings[0]
		last := strings[len(strings)-1]

		for i := 0; i < len(first); i++ {
			if !end && string(last[i]) == string(first[i]) {
				longestPrefix += string(last[i])
			} else {
				end = true
			}
		}
	}

	return longestPrefix
}
