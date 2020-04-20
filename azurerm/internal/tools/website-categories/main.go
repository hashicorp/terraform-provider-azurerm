package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func main() {
	filePath := flag.String("path", "", "The relative path to the output file")
	showHelp := flag.Bool("help", false, "Display this message")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	if err := run(*filePath); err != nil {
		panic(err)
	}
}

func run(outputFileName string) error {
	websiteCategories := make([]string, 0)

	// get a distinct list
	for _, service := range provider.SupportedServices() {
		for _, category := range service.WebsiteCategories() {
			if contains(websiteCategories, category) {
				continue
			}

			websiteCategories = append(websiteCategories, category)
		}
	}

	// sort them
	sort.Strings(websiteCategories)

	outputPath, err := filepath.Abs(outputFileName)
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

	// dump them to a file
	for _, category := range websiteCategories {
		_, _ = file.WriteString(fmt.Sprintf("%s\n", category))
	}

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
