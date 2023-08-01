package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"golang.org/x/tools/go/ast/astutil"
)

var logger hclog.Logger

func main() {
	logger = hclog.New(hclog.DefaultOptions)
	if os.Getenv("DEBUG") != "" {
		logger.SetLevel(hclog.Debug)
	}

	f := flag.NewFlagSet("update-api-version", flag.ExitOnError)
	serviceName := f.String("service", "", "-service=compute")
	oldApiVersion := f.String("old-api-version", "", "-old-api-version=2019-01-01")
	newApiVersion := f.String("new-api-version", "", "-new-api-version=2023-06-01")
	if len(os.Args) == 1 { // 0 is the app name
		log.Fatalf("expected multiple arguments but didn't get any")
	}
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("parsing arguments: %+v", err)
	}
	if serviceName == nil || *serviceName == "" {
		log.Fatalf("missing `-service`")
	}
	if oldApiVersion == nil || *oldApiVersion == "" {
		log.Fatalf("missing `-old-api-version`")
	}
	if newApiVersion == nil || *newApiVersion == "" {
		log.Fatalf("missing `-new-api-version`")
	}

	workingDirectory := fmt.Sprintf("../../services/%s", *serviceName)
	if err := run(*serviceName, *oldApiVersion, *newApiVersion, workingDirectory); err != nil {
		log.Fatalf("error: %+v", err)
	}
}

func run(serviceName string, oldApiVersion string, newApiVersion string, workingDirectory string) error {
	logger.Debug(fmt.Sprintf("Updating Imports in the top-level directory %q..", workingDirectory))
	if err := updateImportsWithinDirectory(serviceName, oldApiVersion, newApiVersion, workingDirectory); err != nil {
		return fmt.Errorf("updating the imports within the top-level directory %q: %+v", workingDirectory, err)
	}

	logger.Debug(fmt.Sprintf("Updating Imports in the directories within directory %q..", workingDirectory))
	entries, err := os.ReadDir(workingDirectory)
	if err != nil {
		return fmt.Errorf("opening the working directory at %q: %+v", workingDirectory, err)
	}
	directories := make([]string, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("retrieving information for %q: %+v", entry.Name(), err)
		}
		if info.IsDir() {
			directories = append(directories, info.Name())
		}
	}
	for _, directory := range directories {
		path := filepath.Join(workingDirectory, directory)
		logger.Debug(fmt.Sprintf("Updating Imports within the nested directory %q..", path))
		if err := updateImportsWithinDirectory(serviceName, oldApiVersion, newApiVersion, path); err != nil {
			return fmt.Errorf("updating the imports within %q: %+v", path, err)
		}
	}

	return nil
}

func updateImportsWithinDirectory(serviceName string, oldApiVersion string, newApiVersion string, workingDirectory string) error {
	fileSet := token.NewFileSet()
	files, err := parser.ParseDir(fileSet, workingDirectory, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing files within %q: %+v", workingDirectory, err)
	}
	for pkgName, pkg := range files {
		fmt.Printf("Package %q", pkgName)
		for fileName, file := range pkg.Files {
			logger.Info(fmt.Sprintf("Updating imports for File %q..", fileName))
			updateImportsForFile(fileSet, file, serviceName, oldApiVersion, newApiVersion)

			var buf bytes.Buffer
			if err = format.Node(&buf, fileSet, file); err != nil {
				return fmt.Errorf("error formatting new code: %w", err)
			}
			_ = os.WriteFile(fileName, buf.Bytes(), 0644)
		}
	}
	return nil
}

func updateImportsForFile(fileSet *token.FileSet, file *ast.File, serviceName string, oldApiVersion string, newApiVersion string) {
	importLineForPreviousApiVersion := fmt.Sprintf("github.com/hashicorp/go-azure-sdk/resource-manager/%s/%s", serviceName, oldApiVersion)
	importLineForNewApiVersion := fmt.Sprintf("github.com/hashicorp/go-azure-sdk/resource-manager/%s/%s", serviceName, newApiVersion)

	existingImports := astutil.Imports(fileSet, file)
	for _, val := range existingImports {
		for _, item := range val {
			logger.Debug(fmt.Sprintf("Processing Import %q", item.Path.Value))
			existingImportLine := item.Path.Value
			if !strings.Contains(existingImportLine, importLineForPreviousApiVersion) {
				continue
			}

			updatedImportLine := strings.Replace(existingImportLine, importLineForPreviousApiVersion, importLineForNewApiVersion, 1)
			logger.Debug(fmt.Sprintf("Updating Import URI from %q to %q", existingImportLine, updatedImportLine))
			item.Path.Value = updatedImportLine

			// if we're importing the meta client (e.g. the api version directly) then we also need to update the alias
			importsMetaClient := existingImportLine == importLineForPreviousApiVersion
			if importsMetaClient && item.Name != nil {
				if existingAlias := item.Name.Name; existingAlias != "" {
					updatedAlias := strings.ToLower(fmt.Sprintf("%s_%s", serviceName, strings.ReplaceAll(newApiVersion, "-", "_")))
					logger.Debug(fmt.Sprintf("Updating Import Alias from %q to %q", existingAlias, updatedAlias))
					item.Name.Name = updatedAlias
				}
			}

			// finally, remove any comments which will be stragglers/lintignores which shouldn't be present
			if item.Comment != nil {
				item.Comment.List = []*ast.Comment{}
			}
		}
	}
}
