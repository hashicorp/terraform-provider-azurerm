// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package data

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

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

	SchemaProperties *Properties

	Document *markdown.Document // resource document
	// These are separated because when it comes to the docs, we may encounter duplicate properties where one is in attributes and another in arguments
	// E.g. `identity` blocks, expect in both args and attrs, but the nested fields should be different
	DocumentArguments  *Properties
	DocumentAttributes *Properties

	Errors []error // errors found in this resource
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

func GetAllTerraformNodeData(fs afero.Fs, providerDir string, serviceName string, resourceName string) []*TerraformNodeData {
	result := make([]*TerraformNodeData, 0)

	pkgData := loadPackages(providerDir)

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

		service.APIsByResource = findAPIsForTypedResources(*pkgData, service)

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

			rd.populateAdditionalFields(fs)

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

			rd.populateAdditionalFields(fs)

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

		service.APIsByResource = findAPIsForUntypedResources(*pkgData, service)

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

			rd.populateAdditionalFields(fs)

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

			rd.populateAdditionalFields(fs)

			result = append(result, rd)
		}
	}

	// TODO: Framework resources
	// for _, s := range provider.SupportedFrameworkServices() {
	//
	// }

	return result
}

func (rd *TerraformNodeData) populateAdditionalFields(fs afero.Fs) {
	rd.populateAPIData()
	rd.populateTimeouts()
	rd.populateDocumentData(fs)
	rd.populateDocumentProperties()
	rd.populateSchemaProperties()
}

func (rd *TerraformNodeData) populateAPIData() {
	if v, ok := rd.Service.APIsByResource[rd.Path]; ok {
		rd.APIs = v
	}
}

func (rd *TerraformNodeData) populateDocumentData(fs afero.Fs) {
	rd.Document.Exists = util.FileExists(fs, rd.Document.Path)

	if rd.Document.Exists {
		if err := rd.Document.Parse(fs); err != nil {
			rd.Errors = append(rd.Errors, fmt.Errorf("failed to parse documentation: %+v", err)) // Output error instead?
		}
	}
}

func (rd *TerraformNodeData) populateSchemaProperties() {
	rd.SchemaProperties = NewProperties()

	populateAllSchemaProperties(rd.SchemaProperties, rd.Resource)
}

func (rd *TerraformNodeData) populateDocumentProperties() {
	var argumentsSection, attributesSection *markdown.Section

	for _, s := range rd.Document.Sections {
		switch s.(type) {
		case *markdown.ArgumentsSection:
			argumentsSection = &s
		case *markdown.AttributesSection:
			attributesSection = &s
		}
	}

	// TODO: move this somewhere else?
	propertyRegex := regexp.MustCompile(`^[-*]\s*\x60([a-z0-9_]*)\x60`)
	blockPropertyRegex := regexp.MustCompile(`^(?i).*[ \t]*\x60(\w*)\x60[ \t]*blocks?.*(?:below|above)`)
	blockSectionRegex := regexp.MustCompile(`^(?i)[\w \t]*\x60(\w*)\x60[ \t]*block[ \t]*(?:supports|exports|contains)`)

	// TODO: consider reworking this? e.g. Gather all blocks then try to map them afterwards
	for _, section := range []*markdown.Section{argumentsSection, attributesSection} {
		if section == nil {
			continue
		}

		// TODO: do differently ot avoid having to deref
		section := *section

		props := NewProperties()

		var lastProperty *Property
		var lastBlockArgumentName string
		var parsingBlockSection bool

		parentlessProps := NewProperties() // tracks blocks for which we've not found a parent yet, e.g. in cases where a block is defined before the parent arg
		fmt.Printf("Processing %s `%s`\n", rd.Type, rd.Name)
		for _, line := range section.GetContent() {
			switch {
			case blockSectionRegex.MatchString(line):
				parsingBlockSection = true

				matches := blockSectionRegex.FindStringSubmatch(line)
				if len(matches) != 2 {
					// TODO: err/debug?
					continue
				}
				lastBlockArgumentName = matches[1]
			case propertyRegex.MatchString(line):
				matches := propertyRegex.FindStringSubmatch(line)
				if len(matches) != 2 {
					// TODO: debug log msg
					continue
				}

				name := matches[1]
				// TODO: func
				isBlock := blockPropertyRegex.MatchString(line)

				prop := Property{
					Name:  name,
					Count: 1,
					Block: isBlock,
				}

				if parsingBlockSection {
					addPropToParent(lastBlockArgumentName, props, &prop, 0)
					continue
				}

				if props == nil { // shouldn't be possible but just in case
					panic("should probably add some error handling ¯\\_(ツ)_/¯")
				}

				if existingArg, ok := props.Objects[name]; ok {
					// if already encountered arg, it may be a duplicate, increment count and let a rule deal with it
					// TODO: when skipping, we may end up with unintentional additionallines on a prop, e.g. an extra ""
					// how to handle?
					existingArg.Count++
					continue
				}

				if isBlock {
					prop.Nested = NewProperties()
				}

				props.Names = append(props.Names, name)
				props.Objects[name] = &prop
				lastProperty = &prop
			case line == "---": // TODO: Is there a better way to track end of block?
				parsingBlockSection = false
				// lastParentArgument = nil
				// lastProperty = nil
			default:
				// default to appending unmatched lines to additional for arg
				if lastProperty != nil { // do we care about excluding `---`?
					lastProperty.AdditionalLines = append(lastProperty.AdditionalLines, line)
				}
			}
		}

		// For any props that weren't matched yet, try to match
		// TODO: do we need to do this multiple times? confirm ¯\_(ツ)_/¯
		lastProperty = nil
		for name, props2 := range parentlessProps.Objects {
			addPropToParent(name, props, props2, 0)
		}

		switch section.(type) {
		case *markdown.ArgumentsSection:
			rd.DocumentArguments = props
		case *markdown.AttributesSection:
			rd.DocumentAttributes = props
		}
	}
}

// TODO: fix stackoverflow panic better, ideally we don't exit based on an arbitrarily set level int
// happening (in one instance) because the `azurerm_dynatrace_monitor` data source has a block named `plan` and within that block is a property named `plan`
func addPropToParent(parentName string, props *Properties, propToAdd *Property, level int) {
	// hacky way to avoid stackoverflow error, TODO: find a better way
	// 10 because having that many nested blocks is absurd and not expected
	if level > 10 {
		return
	}

	if props == nil {
		// TODO err
		return
	}

	if prop, ok := props.Objects[parentName]; ok {
		if prop.Nested == nil {
			prop.Nested = NewProperties() // init if nil
		}

		prop.Nested.Names = append(prop.Nested.Names, propToAdd.Name)
		prop.Nested.Objects[propToAdd.Name] = propToAdd
		// we don't break here, there may be properties in a block that belong to multiple parents
		// e.g. azurerm_function_app_slot.site_config.ip_restriction.headers.* && azurerm_function_app_slot.site_config.scm_ip_restriction.headers.*
	}

	for _, prop := range props.Objects {
		// TODO: consider adding the condition back in, currently unreliable due to the is block regex not matching all block prop definitions
		// A nested argument can only exist in other blocks, skip everything else
		// if prop.Block {
		addPropToParent(parentName, prop.Nested, propToAdd, level+1)
		//}
	}

	return
}

func populateAllSchemaProperties(properties *Properties, resource *schema.Resource) {
	// TODO: rename prop schema
	for name, property := range resource.Schema {
		// TODO guard against propSchema == nil? Realistically shouldn't be possible
		properties.Names = append(properties.Names, name) // This isn't really needed for schema, but we'll leave it in case it's useful later
		properties.Objects[name] = &Property{
			Name:        name,
			Type:        strings.TrimPrefix(property.Type.String(), "Type"),
			Description: property.Description,
			Required:    property.Required,
			Optional:    property.Optional,
			Computed:    property.Computed,
			ForceNew:    property.ForceNew,
			Deprecated:  property.Deprecated != "",
			// PossibleValues:  nil, // TODO
			// DefaultValue:    nil, // TODO
		}

		if r, ok := property.Elem.(*schema.Resource); ok {
			property := properties.Objects[name]

			property.Block = true
			// Expect nested, so init
			// TODO: do we want to check this doesn't override existing? shouldnt be possible but just in case?
			property.Nested = NewProperties()

			populateAllSchemaProperties(property.Nested, r)
		}

		if r, ok := property.Elem.(*schema.Schema); ok {
			properties.Objects[name].NestedType = strings.TrimPrefix(r.Type.String(), "Type")
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
