package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/api-version-lint/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/api-version-lint/version"
)

const (
	EXCEPTIONS_FILE            = "internal/tools/api-version-lint/exceptions.yml"
	HISTORICAL_EXCEPTIONS_FILE = "internal/tools/api-version-lint/historical-exceptions.yml"
	VENDOR_DIR                 = "vendor"
)

func main() {
	historicalExceptions, err := version.ParseHistoricalExceptions(filepath.FromSlash(HISTORICAL_EXCEPTIONS_FILE))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse historical exceptions file %s:\n\n%s\n\n", HISTORICAL_EXCEPTIONS_FILE, err)
		printErrorFooterAndExit()
	}
	exceptions, err := version.ParseExceptions(filepath.FromSlash(EXCEPTIONS_FILE))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse exceptions file %s:\n\n%s\n\n", EXCEPTIONS_FILE, err)
		printErrorFooterAndExit()
	}

	previewVersions := map[version.Version]bool{}
	for _, sdkType := range sdk.SdkTypes {
		if err := populatePreviewVersions(sdkType, previewVersions); err != nil {
			fmt.Fprintf(os.Stderr, "failed to populate preview versions for SDK type %s:\n\n%s\n\n", sdkType.Module, err)
			printErrorFooterAndExit()
		}
	}

	unusedExceptions := []string{}
	for _, exception := range historicalExceptions {
		if !previewVersions[exception] {
			unusedExceptions = append(unusedExceptions, fmt.Sprintf("module: %s, service: %s, version: %s", exception.Module, exception.Service, exception.Version))
		} else {
			delete(previewVersions, exception)
		}
	}
	failIfAnyUnusuedExceptions(unusedExceptions, HISTORICAL_EXCEPTIONS_FILE)

	for _, exception := range exceptions {
		if !previewVersions[exception] {
			unusedExceptions = append(unusedExceptions, fmt.Sprintf("module: %s, service: %s, version: %s", exception.Module, exception.Service, exception.Version))
		} else {
			delete(previewVersions, exception)
		}
	}
	failIfAnyUnusuedExceptions(unusedExceptions, EXCEPTIONS_FILE)

	if len(previewVersions) > 0 {
		invalidPreviewVersions := []string{}
		for svcVer := range previewVersions {
			invalidPreviewVersions = append(invalidPreviewVersions, fmt.Sprintf("module: %s, service: %s, version: %s", svcVer.Module, svcVer.Service, svcVer.Version))
		}
		sort.Strings(invalidPreviewVersions)

		fmt.Fprintf(os.Stderr, "❌ Invalid use of preview SDK API versions detected in `vendor` folder:\n\n")
		for _, v := range invalidPreviewVersions {
			fmt.Fprintf(os.Stderr, "%s\n", v)
		}
		fmt.Fprintf(os.Stderr, `
Preview versions are prone to breaking changes resulting in very bad user experience, 
please use stable version instead.

`)
		printErrorFooterAndExit()
	}

	fmt.Printf("✅ No invalid use of preview SDK API versions detected in %s folder.\n", VENDOR_DIR)
}

func populatePreviewVersions(sdkType sdk.SdkType, previewServiceVersions map[version.Version]bool) error {
	return filepath.WalkDir(filepath.FromSlash(VENDOR_DIR+"/"+sdkType.Module), func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dir.IsDir() {
			matches := sdkType.ServiceAndVersionRegex.FindStringSubmatch(filepath.ToSlash(path))
			if len(matches) == 3 {
				previewServiceVersions[version.Version{
					Module:  sdkType.Module,
					Service: strings.ToLower(matches[1]),
					Version: strings.ToLower(matches[2]),
				}] = true
			}
		}
		return nil
	})
}

func printErrorFooterAndExit() {
	fmt.Fprintf(os.Stderr, `More information: https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-api-version.md

To rerun this check locally, use: go run internal/tools/api-version-lint/main.go
`)
	os.Exit(1)
}

func failIfAnyUnusuedExceptions(unusedExceptions []string, exceptionsFile string) {
	if len(unusedExceptions) > 0 {
		fmt.Fprintf(os.Stderr, "❌ Unused exceptions detected in `%s` file, remove these entries:\n\n", exceptionsFile)
		for _, unusedException := range unusedExceptions {
			fmt.Fprintln(os.Stderr, unusedException)
		}
		fmt.Fprintf(os.Stderr, "\n")
		printErrorFooterAndExit()
	}
}
