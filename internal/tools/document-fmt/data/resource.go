// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
	"github.com/spf13/afero"
)

type ResourceType string

const (
	ResourceTypeData      ResourceType = "Data Source"
	ResourceTypeEphemeral ResourceType = "Ephemeral Resource"
	ResourceTypeResource  ResourceType = "Resource"
)

var (
	resourceFilePathPattern    = "%s/%s_%s.go"
	resourceFileGenPathPattern = "%s/%s_%s_gen.go"

	ResourceTypeToFileSuffix = map[ResourceType]string{
		ResourceTypeData:      "data_source",
		ResourceTypeEphemeral: "ephemeral",
		ResourceTypeResource:  "resource",
	}

	ResourceTypeToDocumentationSubPath = map[ResourceType]string{
		ResourceTypeData:      "d",
		ResourceTypeEphemeral: "ephemeral-resources",
		ResourceTypeResource:  "r",
	}
)

func (r ResourceType) String() string {
	return string(r)
}

func expectedResourceCodePath(fs afero.Fs, pattern string, name string, service Service, resourceType ResourceType) string {
	defaultPath := filepath.FromSlash(fmt.Sprintf(pattern, service.Path, name, ResourceTypeToFileSuffix[resourceType]))
	if fs == nil || util.FileExists(fs, defaultPath) {
		return defaultPath
	}

	suffix := ResourceTypeToFileSuffix[resourceType]
	if strings.Contains(pattern, "_gen.go") {
		suffix += "_gen"
	}

	cleanServiceName := strings.ToLower(strings.ReplaceAll(service.Name, " ", ""))
	cleanServiceNameUnderscore := strings.ToLower(strings.ReplaceAll(service.Name, " ", "_"))
	serviceBaseDir := filepath.Base(service.Path)
	shortNameWithoutService := strings.TrimPrefix(name, cleanServiceName+"_")
	shortNameWithoutService = strings.TrimPrefix(shortNameWithoutService, cleanServiceNameUnderscore+"_")
	shortNameWithoutService = strings.TrimPrefix(shortNameWithoutService, serviceBaseDir+"_")

	candidates := []string{
		filepath.Join(service.Path, shortNameWithoutService, fmt.Sprintf("%s.go", suffix)),
		filepath.Join(service.Path, shortNameWithoutService, fmt.Sprintf("%s_%s.go", shortNameWithoutService, suffix)),
		filepath.Join(service.Path, name, fmt.Sprintf("%s.go", suffix)),
		filepath.Join(service.Path, name, fmt.Sprintf("%s_%s.go", name, suffix)),
	}

	if idx := strings.LastIndex(name, "_"); idx != -1 {
		lastPart := name[idx+1:]
		candidates = append(candidates,
			filepath.Join(service.Path, lastPart, fmt.Sprintf("%s.go", suffix)),
			filepath.Join(service.Path, lastPart, fmt.Sprintf("%s_%s.go", lastPart, suffix)),
		)
	}

	for _, candidate := range candidates {
		if util.FileExists(fs, candidate) {
			return candidate
		}
	}

	return defaultPath
}
