// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/container-apps"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Container Apps",
	}
}

func (r Registration) Name() string {
	return "Container Apps"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ContainerAppDataSource{},
		ContainerAppEnvironmentDataSource{},
		ContainerAppEnvironmentCertificateDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ContainerAppEnvironmentCertificateResource{},
		ContainerAppEnvironmentCustomDomainResource{},
		ContainerAppEnvironmentDaprComponentResource{},
		ContainerAppEnvironmentResource{},
		ContainerAppEnvironmentStorageResource{},
		ContainerAppResource{},
		ContainerAppCustomDomainResource{},
		ContainerAppJobResource{},
	}
}
