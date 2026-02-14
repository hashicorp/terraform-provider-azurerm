// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helper

import "strings"

// ShouldSkipPackageForResourceAnalysis returns true if the package should be skipped
// during resource/data source analysis (e.g., test, migration, client, validate packages)
func ShouldSkipPackageForResourceAnalysis(pkgPath string) bool {
	skipPackages := []string{
		"_test",
		"/migration",
		"/client",
		"/validate",
		"/test-data",
		"/parse",
		"/models",
	}

	for _, skip := range skipPackages {
		if strings.Contains(pkgPath, skip) {
			return true
		}
	}
	return false
}

// IsCachePath returns true if the path is a build/cache path that should be skipped
func IsCachePath(filePath string) bool {
	cachePatterns := []string{
		"go-build",
		"AppData",
		".test",
	}

	for _, pattern := range cachePatterns {
		if strings.Contains(filePath, pattern) {
			return true
		}
	}
	return false
}

// IsResourceOrDataSourceFile returns true if the file is a resource or data source file
// (i.e., ends with _resource.go or _data_source.go)
func IsResourceOrDataSourceFile(filename string) bool {
	resourceFileSuffix := []string{
		"_resource.go",
		"_data_source.go",
	}

	for _, suffix := range resourceFileSuffix {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}
