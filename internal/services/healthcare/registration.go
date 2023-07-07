// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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
		"azurerm_healthcare_service":         dataSourceHealthcareService(),
		"azurerm_healthcare_workspace":       dataSourceHealthcareWorkspace(),
		"azurerm_healthcare_dicom_service":   dataSourceHealthcareDicomService(),
		"azurerm_healthcare_fhir_service":    dataSourceHealthcareApisFhirService(),
		"azurerm_healthcare_medtech_service": dataSourceHealthcareIotConnector(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_healthcare_service":                          resourceHealthcareService(),
		"azurerm_healthcare_workspace":                        resourceHealthcareApisWorkspace(),
		"azurerm_healthcare_dicom_service":                    resourceHealthcareApisDicomService(),
		"azurerm_healthcare_fhir_service":                     resourceHealthcareApisFhirService(),
		"azurerm_healthcare_medtech_service":                  resourceHealthcareApisMedTechService(),
		"azurerm_healthcare_medtech_service_fhir_destination": resourceHealthcareApisMedTechServiceFhirDestination(),
	}
}
