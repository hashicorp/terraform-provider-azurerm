// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
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
	// new clusters cannot be created - so there's no point trying to run the tests
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
- changed-files:
  - any-glob-to-any-file:
    - go.mod
    - go.sum
    - vendor/**/*
documentation:
- changed-files:
  - any-glob-to-any-file: 
    - website/**/*
tooling:
- changed-files:
  - any-glob-to-any-file: 
    - internal/tools/**/*
state-migration:
- changed-files:
  - any-glob-to-any-file: 
    - internal/services/**/migration/**/*
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
			"- changed-files:",
			"  - any-glob-to-any-file:",
		}

		for _, pkg := range pkgs {
			out = append(out, fmt.Sprintf("    - internal/services/%[1]s/**/*", pkg))
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
	template := `// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// NOTE: this is Generated from the Service Definitions - manual changes will be lost
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
v/1.x (legacy):
  - '### AzureRM Provider Version\s+(|azurerm |AzureRM )(|v|V)1\.\d+'
v/2.x (legacy):
  - '### AzureRM Provider Version\s+(|azurerm |AzureRM )(|v|V)2\.\d+'
v/3.x:
  - '### AzureRM Provider Version\s+(|azurerm |AzureRM )(|v|V)3\.\d+'
v/4.x:
  - '### AzureRM Provider Version\s+(|azurerm |AzureRM )(|v|V)4\.\d+'
`

const azurerm = "azurerm_"

type githubIssueLabelsGenerator struct{}

func (g githubIssueLabelsGenerator) outputPath(rootDirectory string) string {
	return fmt.Sprintf("%s/.github/labeler-issue-triage.yml", rootDirectory)
}

type Prefix struct {
	Names        []string
	CommonPrefix string
}

func (githubIssueLabelsGenerator) run(outputFileName string, _ map[string]struct{}) error {
	labelToNames := make(map[string][]string)
	label := ""

	for _, service := range provider.SupportedTypedServices() {

		v, ok := service.(sdk.TypedServiceRegistrationWithAGitHubLabel)
		// keep a record of resources/datasources that don't have labels so they can be used to check that prefixes generated later don't match resources from those services
		label = ""
		if ok {
			label = v.AssociatedGitHubLabel()
		}

		var names []string
		for _, resource := range service.Resources() {
			names = append(names, resource.ResourceType())
		}

		for _, ds := range service.DataSources() {
			names = append(names, ds.ResourceType())
		}

		labelToNames = appendToSliceWithinMap(labelToNames, names, label)

	}
	for _, service := range provider.SupportedUntypedServices() {
		v, ok := service.(sdk.UntypedServiceRegistrationWithAGitHubLabel)
		service.SupportedResources()

		// keep a record of resources/datasources that don't have labels so they can be used to check that prefixes generated later don't match resources from those services
		label = ""
		if ok {
			label = v.AssociatedGitHubLabel()
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

		names = removeDuplicateNames(names)
		labelToNames = appendToSliceWithinMap(labelToNames, names, label)

	}

	sortedLabels := make([]string, 0)
	for k := range labelToNames {
		sortedLabels = append(sortedLabels, k)
	}
	sort.Strings(sortedLabels)

	output := strings.TrimSpace(githubIssueLabelsTemplate)

	labelToPrefixes := make(map[string][]Prefix)

	// loop through all labels and get a list of prefixes that match each label. And for each prefix, record which resource/datasource names it is derived from - we need to retain these in case there are duplicate prefixes matching resources with a different label
	for _, labelName := range sortedLabels {

		longestPrefix := longestCommonPrefix(labelToNames[labelName])
		var prefixGroups []Prefix
		// If there is no common prefix for a service, separate it into groups using the next segment of the name (azurerm_xxx) and add multiple possible prefixes for the label
		// For example, under "service/signalr" we separate names into groups of 2 prefixes.
		// The prefix azurerm_signalr matches resources `azurerm_signalr_shared_private_link_resource`, `azurerm_signalr_service`, `azurerm_signalr_service_custom_domain` etc.
		// And web_pubsub matches `azurerm_web_pubsub_hub`, `azurerm_web_pubsub_network_acl` etc.
		// But both share the "service/signalr" label.
		if longestPrefix == azurerm {
			prefixGroups = getPrefixesForNames(labelToNames[labelName])
		} else {
			prefixGroups = []Prefix{
				{
					Names:        labelToNames[labelName],
					CommonPrefix: longestPrefix,
				},
			}
		}
		labelToPrefixes[labelName] = prefixGroups

	}

	// loop though again, this time compiling prefixes into a regex for each label and separating out duplicates
	for _, labelName := range sortedLabels {

		if labelName == "" {
			continue
		}

		out := []string{
			fmt.Sprintf("%[1]s:", labelName),
		}
		prefixes := make([]string, 0)

		for _, prefix := range labelToPrefixes[labelName] {
			// if a prefix matches another prefix, use the whole name for each resource/ds that matches that prefix in the regex
			if prefixHasMatch(labelName, prefix, labelToPrefixes) {

				for _, name := range prefix.Names {
					prefixes = append(prefixes, strings.TrimPrefix(name+"\\W+", azurerm))
				}
			} else {
				prefixes = append(prefixes, strings.TrimPrefix(prefix.CommonPrefix, azurerm))
			}
		}

		if len(prefixes) > 0 {
			if len(prefixes) > 1 {
				out = append(out, fmt.Sprintf("  - '### (|New or )Affected Resource\\(s\\)\\/Data Source\\(s\\)((.|\\n)*)azurerm_(%s)((.|\\n)*)###'", strings.Join(prefixes, "|")))
			}
			if len(prefixes) == 1 {
				out = append(out, fmt.Sprintf("  - '### (|New or )Affected Resource\\(s\\)\\/Data Source\\(s\\)((.|\\n)*)azurerm_%s((.|\\n)*)###'", prefixes[0]))
			}
			out = append(out, "")
			output += fmt.Sprintf("\n%s", strings.Join(out, "\n"))
		}
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

func longestCommonPrefix(names []string) string {
	longestPrefix := ""
	end := false

	if len(names) > 0 {

		sort.Strings(names)
		first := names[0]
		last := names[len(names)-1]

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

func commonPrefixGroups(names []string) [][]string {
	var prefixGroups [][]string

	if len(names) > 0 {
		sort.Strings(names)

		// get the prefix of the first name in the list matching azurerm_xxxx or azurerm_xxx_* and add it to the first group
		re := regexp.MustCompile("azurerm_[a-zA-Z0-9]+($|_)")
		prefix := ""
		if match := re.FindStringSubmatch(names[0]); match != nil {
			prefix = strings.TrimSuffix(match[0], "_")
		}
		group := []string{names[0]}

		// loop through the rest of the names beginning with the second in the list
		// adding names to the group when their prefix matches the prefix of the previous name
		for i := 1; i < len(names); i++ {
			currentPrefix := ""
			if match := re.FindStringSubmatch(names[i]); match != nil {
				currentPrefix = strings.TrimSuffix(match[0], "_")
			}
			if currentPrefix == prefix {
				group = append(group, names[i])
			} else {
				// append group and start a new prefix group if the current prefix does not match the prefix from the previous group
				prefixGroups = append(prefixGroups, group)
				group = []string{names[i]}
				prefix = currentPrefix
			}
		}
		prefixGroups = append(prefixGroups, group)
	}
	return prefixGroups
}

func prefixHasMatch(labelToCheck string, prefixToCheck Prefix, labelToPrefixes map[string][]Prefix) bool {
	for label, allPrefixes := range labelToPrefixes {
		if label == labelToCheck {
			continue
		}
		for _, prefix := range allPrefixes {
			if strings.Contains(prefix.CommonPrefix, prefixToCheck.CommonPrefix) {
				return true
			}
		}
	}

	return false
}

func getPrefixesForNames(names []string) []Prefix {

	var prefixes []Prefix
	groupedNames := commonPrefixGroups(names)

	for _, group := range groupedNames {
		prefix := Prefix{
			Names:        group,
			CommonPrefix: longestCommonPrefix(group),
		}
		prefixes = append(prefixes, prefix)
	}

	return prefixes
}

func removeDuplicateNames(names []string) []string {
	keys := make(map[string]bool)
	uniqueNames := make([]string, 0)

	for _, name := range names {
		if _, value := keys[name]; !value {
			keys[name] = true
			uniqueNames = append(uniqueNames, name)
		}
	}
	return uniqueNames
}
