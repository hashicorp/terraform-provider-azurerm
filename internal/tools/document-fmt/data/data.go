// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package data

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	gomonkey "github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/mdparser"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// patchPossibleValuesFn patches validation.StringInSlice to capture possible values
func patchPossibleValuesFn() {
	gomonkey.ApplyFunc(validation.StringInSlice,
		func(valid []string, ignoreCase bool) schema.SchemaValidateFunc { //nolint:staticcheck
			return func(i interface{}, k string) (warnings []string, errors []error) {
				var res []string // must have a copy
				res = append(res, valid...)
				return res, nil
			}
		})
}

func init() {
	patchPossibleValuesFn()
}

// TerraformNodeData contains all data for a resource that may be required for validation/scaffolding
// TODO: Provider type (struct) containing map of Service type, which should each contain a slice of `TerraformNodeData` structs
type TerraformNodeData struct {
	Name         string // resource name
	ShortName    string // resource name minus provider prefix
	ProviderName string // provider name

	Service  Service          // resource's service package
	Type     ResourceType     // resource type
	Path     string           // resource expected code file path
	Resource *schema.Resource // sdk resource

	APIs     []API     // APIs used by this resource -- best effort, may not be populated
	Timeouts []Timeout // Timeouts from *schema.Resource

	SchemaProperties *models.SchemaProperties

	Document           *markdown.Document // resource document
	DocumentArguments  *models.DocumentProperties
	DocumentAttributes *models.DocumentProperties

	ShouldNormalize bool // whether to normalize document before parsing
	Errors          []error
}

func newTerraformNodeData(fs afero.Fs, providerDir string, service Service, name string, resourceType ResourceType, source any) (*TerraformNodeData, error) {
	providerName, _, _ := strings.Cut(name, "_")

	result := TerraformNodeData{
		Name:         name,
		ShortName:    strings.TrimPrefix(name, fmt.Sprintf("%s_", providerName)),
		ProviderName: providerName,
		Service:      service,
		Type:         resourceType,
	}
	result.Path = expectedResourceCodePath(resourceFilePathPattern, result.ShortName, service, resourceType)

	// skip if generated resource
	if util.FileExists(fs, expectedResourceCodePath(resourceFileGenPathPattern, result.ShortName, service, resourceType)) {
		return nil, fmt.Errorf("skipping generated resource") // TODO debug msg, no error?
	}

	result.Document = markdown.NewDocument(expectedDocumentationPath(providerDir, result.ShortName, result.Type))

	switch r := source.(type) {
	case sdk.Resource:
		w := sdk.NewResourceWrapper(r)
		wr, err := w.Resource()
		if err != nil {
			return nil, fmt.Errorf("wrapping resource: %+v", err)
		}
		result.Resource = wr
	case sdk.DataSource:
		w := sdk.NewDataSourceWrapper(r)
		wr, err := w.DataSource()
		if err != nil {
			return nil, fmt.Errorf("wrapping data source: %+v", err)
		}
		result.Resource = wr
	case *schema.Resource:
		result.Resource = r
	default:
		return nil, fmt.Errorf("unexpected type `%T` for resource `%s`", r, result.ShortName)
	}

	return &result, nil
}

func GetAllTerraformNodeData(fs afero.Fs, providerDir string, serviceName string, resourceName string, shouldNormalizeMd bool, shouldLoadPackages bool) []*TerraformNodeData {
	result := make([]*TerraformNodeData, 0)

	// Only load packages if needed by the selected rules
	var pkgData *packageData
	if shouldLoadPackages {
		log.WithField("reason", "required by selected rules").Info("loading packages for API analysis")
		pkgData = loadPackages(providerDir)
	} else {
		log.Info("skipping package loading - no rules require it")
	}

	for _, s := range provider.SupportedTypedServices() {
		service, err := NewService(fs, providerDir, s, s.Name())
		if err != nil {
			log.WithFields(log.Fields{
				"service": s.Name(),
				"error":   err,
			}).Warn("Skipping service...")
			continue
		}

		// TODO Skip based on multiple services?
		if serviceName != "" {
			if service.Name != serviceName {
				continue
			}
		}

		if pkgData != nil {
			service.APIsByResource = findAPIsForTypedResources(*pkgData, service)
		}

		for _, r := range s.DataSources() {
			name := r.ResourceType()

			// TODO Skip based on multiple resources?
			if resourceName != "" {
				if name != resourceName {
					continue
				}
			}

			rd, err := newTerraformNodeData(fs, providerDir, *service, name, ResourceTypeData, r)
			if err != nil {
				log.Error(err)
				continue
			}

			rd.populateAdditionalFields(fs, false)

			result = append(result, rd)
		}

		for _, r := range s.Resources() {
			name := r.ResourceType()

			// TODO Skip based on multiple resources?
			if resourceName != "" {
				if name != resourceName {
					continue
				}
			}

			rd, err := newTerraformNodeData(fs, providerDir, *service, name, ResourceTypeResource, r)
			if err != nil {
				log.Error(err)
				continue
			}

			rd.populateAdditionalFields(fs, shouldNormalizeMd)

			result = append(result, rd)
		}
	}
	for _, s := range provider.SupportedUntypedServices() {
		service, err := NewService(fs, providerDir, s, s.Name())
		if err != nil {
			log.WithFields(log.Fields{
				"service": s.Name(),
				"error":   err,
			}).Warn("Skipping Service")
			continue
		}

		// TODO Skip based on multiple services?
		if serviceName != "" {
			if service.Name != serviceName {
				continue
			}
		}

		if pkgData != nil {
			service.APIsByResource = findAPIsForUntypedResources(*pkgData, service)
		}

		for name, r := range s.SupportedDataSources() {
			rd, err := newTerraformNodeData(fs, providerDir, *service, name, ResourceTypeData, r)
			if err != nil {
				log.Error(err)
				continue
			}

			// TODO Skip based on multiple resources?
			if resourceName != "" {
				if name != resourceName {
					continue
				}
			}

			rd.populateAdditionalFields(fs, false)

			result = append(result, rd)
		}

		for name, r := range s.SupportedResources() {
			rd, err := newTerraformNodeData(fs, providerDir, *service, name, ResourceTypeResource, r)
			if err != nil {
				log.Error(err)
				continue
			}

			// TODO Skip based on multiple resources?
			if resourceName != "" {
				if name != resourceName {
					continue
				}
			}

			rd.populateAdditionalFields(fs, shouldNormalizeMd)

			result = append(result, rd)
		}
	}

	// TODO: Framework resources
	// for _, s := range provider.SupportedFrameworkServices() {
	//
	// }

	return result
}

func (rd *TerraformNodeData) populateAdditionalFields(fs afero.Fs, shouldNormalizeMd bool) {
	rd.populateAPIData()
	rd.populateTimeouts()
	rd.populateDocumentData(fs, shouldNormalizeMd)
	rd.populateDocumentProperties()
	rd.populateSchemaProperties()
}

func (rd *TerraformNodeData) populateAPIData() {
	if v, ok := rd.Service.APIsByResource[rd.Path]; ok {
		rd.APIs = v
	}
}

func (rd *TerraformNodeData) populateDocumentData(fs afero.Fs, shouldNormalizeMd bool) {
	rd.Document.Exists = util.FileExists(fs, rd.Document.Path)

	if rd.Document.Exists {
		if err := rd.Document.Parse(fs, shouldNormalizeMd); err != nil {
			rd.Errors = append(rd.Errors, fmt.Errorf("failed to parse documentation: %+v", err)) // Output error instead?
		}
	}
}

func (rd *TerraformNodeData) populateSchemaProperties() {
	rd.SchemaProperties = models.NewSchemaProperties()

	populateAllSchemaProperties(rd.SchemaProperties, rd.Resource)
}

func (rd *TerraformNodeData) populateDocumentProperties() {
	argumentsSection := rd.Document.GetArgumentsSection()

	if argumentsSection != nil {
		if parsedProps := parseMdArgToProperties(argumentsSection); parsedProps != nil {
			rd.DocumentArguments = parsedProps
		}
	}
}

func parseMdArgToProperties(argSection *markdown.ArgumentsSection) *models.DocumentProperties {
	return mdparser.ParseMarkdownSection(argSection.GetContent())
}

// extractPossibleValues attempts to extract possible values from ValidateFunc using monkey patch
func extractPossibleValues(property *schema.Schema) []string {
	if property.ValidateFunc == nil {
		return nil
	}

	// Check if it is StringInSlice validation
	pc := reflect.ValueOf(property.ValidateFunc).Pointer()
	fn := runtime.FuncForPC(pc).Name()

	// ValidateFunc may direct use sdk v2.StringInSlice then function name will be patchPossibleValuesFn
	if strings.Contains(fn, "patchPossibleValuesFn") || strings.Contains(fn, "StringInSlice") {
		values, _ := property.ValidateFunc(nil, "")
		return values
	}

	return nil
}

// extractDefaultValue extracts the default value from schema
func extractDefaultValue(property *schema.Schema) interface{} {
	if property.Default != nil {
		return property.Default
	}

	if property.DefaultFunc != nil {
		if val, err := property.DefaultFunc(); err == nil && val != nil {
			return val
		}
	}

	return nil
}

func populateAllSchemaProperties(properties *models.SchemaProperties, resource *schema.Resource) {
	for name, property := range resource.Schema {
		properties.Objects[name] = &models.SchemaProperty{
			Name:           name,
			Type:           strings.TrimPrefix(property.Type.String(), "Type"),
			Description:    property.Description,
			Required:       property.Required,
			Optional:       property.Optional,
			Computed:       property.Computed,
			ForceNew:       property.ForceNew,
			Deprecated:     property.Deprecated != "",
			PossibleValues: extractPossibleValues(property),
			DefaultValue:   extractDefaultValue(property),
		}

		// Handle nested resources
		if r, ok := property.Elem.(*schema.Resource); ok {
			schemaProperty := properties.Objects[name]
			schemaProperty.Block = true
			schemaProperty.Nested = models.NewSchemaProperties()
			populateAllSchemaProperties(schemaProperty.Nested, r)
		}

		// Handle nested schemas (e.g., List of String with validation)
		if r, ok := property.Elem.(*schema.Schema); ok {
			properties.Objects[name].NestedType = strings.TrimPrefix(r.Type.String(), "Type")
			// Extract possible values from nested schema if available
			if nestedPossibleValues := extractPossibleValues(r); len(nestedPossibleValues) > 0 {
				properties.Objects[name].PossibleValues = nestedPossibleValues
			}
		}
	}
}

func (rd *TerraformNodeData) populateTimeouts() {
	if t := rd.Resource.Timeouts; t != nil {
		if t.Create != nil {
			rd.Timeouts = append(rd.Timeouts, Timeout{
				Type:     TimeoutTypeCreate,
				Duration: int(t.Create.Minutes()),
				Name:     "<Azure Brand Name>",
			})
		}

		if t.Read != nil {
			rd.Timeouts = append(rd.Timeouts, Timeout{
				Type:     TimeoutTypeRead,
				Duration: int(t.Read.Minutes()),
				Name:     "<Azure Brand Name>",
			})
		}

		if t.Update != nil {
			rd.Timeouts = append(rd.Timeouts, Timeout{
				Type:     TimeoutTypeUpdate,
				Duration: int(t.Update.Minutes()),
				Name:     "<Azure Brand Name>",
			})
		}

		if t.Delete != nil {
			rd.Timeouts = append(rd.Timeouts, Timeout{
				Type:     TimeoutTypeDelete,
				Duration: int(t.Delete.Minutes()),
				Name:     "<Azure Brand Name>",
			})
		}
	}
}
