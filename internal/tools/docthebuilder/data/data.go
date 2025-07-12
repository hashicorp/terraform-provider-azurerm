package data

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/markdown"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// ResourceData contains all data for a resource that may be required for validation/scaffolding
type ResourceData struct {
	Name         string // resource name
	ShortName    string // resource name minus provider prefix
	ProviderName string // provider name

	Service  Service          // resource's service package
	Type     ResourceType     // resource type
	Path     string           // resource expected code file path
	Resource *schema.Resource // sdk resource

	APIs     []API     // APIs used by this resource -- best effort, may not be populated
	Timeouts []Timeout // Timeouts from *schema.Resource

	Document *markdown.Document // resource document

	Errors []error // errors found in this resource
}

func newDataObject(fs afero.Fs, providerDir string, service Service, name string, resourceType ResourceType, source any) (*ResourceData, error) {
	providerName, _, _ := strings.Cut(name, "_")

	result := ResourceData{
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

func GetData(fs afero.Fs, providerDir string, serviceName string, resourceName string) []*ResourceData {
	result := make([]*ResourceData, 0)

	pkgData := loadPackages(providerDir)

	for _, s := range provider.SupportedTypedServices() {
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

		service.APIsByResource = findAPIsForTypedResources(*pkgData, service)

		for _, r := range s.DataSources() {
			name := r.ResourceType()

			// TODO Skip based on multiple resources?
			if resourceName != "" {
				if name != resourceName {
					continue
				}
			}

			rd, err := newDataObject(fs, providerDir, *service, name, ResourceTypeData, r)
			if err != nil {
				log.Error(err)
				continue
			}

			rd.populateAPIData()
			rd.populateTimeouts()
			rd.populateDocumentData(fs)

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

			rd, err := newDataObject(fs, providerDir, *service, name, ResourceTypeResource, r)
			if err != nil {
				log.Error(err)
				continue
			}

			rd.populateAPIData()
			rd.populateTimeouts()
			rd.populateDocumentData(fs)

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
			rd, err := newDataObject(fs, providerDir, *service, name, ResourceTypeData, r)
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

			rd.populateAPIData()
			rd.populateTimeouts()
			rd.populateDocumentData(fs)

			result = append(result, rd)
		}

		for name, r := range s.SupportedResources() {
			rd, err := newDataObject(fs, providerDir, *service, name, ResourceTypeResource, r)
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

			rd.populateAPIData()
			rd.populateTimeouts()
			rd.populateDocumentData(fs)

			result = append(result, rd)
		}
	}

	// TODO: Framework resources
	// for _, s := range provider.SupportedFrameworkServices() {
	//
	// }

	return result
}

func (rd *ResourceData) populateAPIData() {
	if v, ok := rd.Service.APIsByResource[rd.Path]; ok {
		rd.APIs = v
	}
}

func (rd *ResourceData) populateDocumentData(fs afero.Fs) {
	rd.Document.Exists = util.FileExists(fs, rd.Document.Path)

	if rd.Document.Exists {
		if err := rd.Document.Parse(fs); err != nil {
			rd.Errors = append(rd.Errors, fmt.Errorf("failed to parse documentation: %+v", err)) // Output error instead?
		}
	}
}

func (rd *ResourceData) populateTimeouts() {
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
