// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package codesigning

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/trustedsigning"
}

// Name is the name of this Service
func (r Registration) Name() string {
	if !features.FivePointOh() {
		return "Trusted Signing"
	}
	return "Artifacts Signing"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	if !features.FivePointOh() {
		return []string{
			"Trusted Signing",
			"Artifacts Signing",
		}
	}
	return []string{
		"Artifacts Signing",
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	if !features.FivePointOh() {
		return []sdk.DataSource{
			ArtifactsSigningAccountDataSource{},
			TrustedSigningAccountDataSource{},
		}
	}
	return []sdk.DataSource{
		ArtifactsSigningAccountDataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	if !features.FivePointOh() {
		return []sdk.Resource{
			ArtifactsSigningAccountResource{},
			TrustedSigningAccountResource{},
		}
	}
	return []sdk.Resource{
		ArtifactsSigningAccountResource{},
	}
}
