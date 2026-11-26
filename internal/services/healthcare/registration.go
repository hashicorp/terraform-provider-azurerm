// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_dicom_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_fhir_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_medtech_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_medtech_service_fhir_destination"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/healthcare"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Health Care"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Healthcare",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_healthcare_dicom_service":   healthcare_dicom_service.RegisterDataSource(),
		"azurerm_healthcare_fhir_service":    healthcare_fhir_service.RegisterDataSource(),
		"azurerm_healthcare_medtech_service": healthcare_medtech_service.RegisterDataSource(),
		"azurerm_healthcare_service":         healthcare_service.RegisterDataSource(),
		"azurerm_healthcare_workspace":       healthcare_workspace.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_healthcare_dicom_service":                    healthcare_dicom_service.RegisterResource(),
		"azurerm_healthcare_fhir_service":                     healthcare_fhir_service.RegisterResource(),
		"azurerm_healthcare_medtech_service":                  healthcare_medtech_service.RegisterResource(),
		"azurerm_healthcare_medtech_service_fhir_destination": healthcare_medtech_service_fhir_destination.RegisterResource(),
		"azurerm_healthcare_service":                          healthcare_service.RegisterResource(),
		"azurerm_healthcare_workspace":                        healthcare_workspace.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
