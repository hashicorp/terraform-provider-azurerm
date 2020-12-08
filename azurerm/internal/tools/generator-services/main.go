package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func main() {
	filePath := flag.String("path", "", "The relative path to the root directory")
	showHelp := flag.Bool("help", false, "Display this message")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	generators := []generator{
		teamCityServicesListGenerator{},
		websiteCategoriesGenerator{},
	}
	for _, value := range generators {
		outputFile := value.outputPath(*filePath)
		if err := value.run(outputFile); err != nil {
			panic(err)
		}
	}
}

type generator interface {
	outputPath(rootDirectory string) string
	run(outputFileName string) error
}

type teamCityServicesListGenerator struct{}

func (teamCityServicesListGenerator) outputPath(rootDirectory string) string {
	return fmt.Sprintf("%s/.teamcity/components/generated/services.kt", rootDirectory)
}

func (teamCityServicesListGenerator) run(outputFileName string) error {
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

func (websiteCategoriesGenerator) run(outputFileName string) error {
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
