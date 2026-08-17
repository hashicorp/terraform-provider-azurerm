// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elasticsan

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/elasticsan"
}

func (Registration) Name() string {
	return "ElasticSan"
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ElasticSANDataSource{},
		ElasticSANVolumeGroupDataSource{},
		ElasticSANVolumeSnapshotDataSource{},
	}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ElasticSANResource{},
		ElasticSANVolumeGroupResource{},
		ElasticSANVolumeResource{},
	}
}

func (Registration) WebsiteCategories() []string {
	return []string{
		"Elastic SAN",
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
