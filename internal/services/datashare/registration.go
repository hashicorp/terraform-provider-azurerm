// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/data-share"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Data Share"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Data Share",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_share_account":                dataSourceDataShareAccount(),
		"azurerm_data_share":                        dataSourceDataShare(),
		"azurerm_data_share_dataset_blob_storage":   dataSourceDataShareDatasetBlobStorage(),
		"azurerm_data_share_dataset_data_lake_gen2": dataSourceDataShareDatasetDataLakeGen2(),
		"azurerm_data_share_dataset_kusto_cluster":  dataSourceDataShareDatasetKustoCluster(),
		"azurerm_data_share_dataset_kusto_database": dataSourceDataShareDatasetKustoDatabase(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_share_account":                resourceDataShareAccount(),
		"azurerm_data_share":                        resourceDataShare(),
		"azurerm_data_share_dataset_blob_storage":   resourceDataShareDataSetBlobStorage(),
		"azurerm_data_share_dataset_data_lake_gen2": resourceDataShareDataSetDataLakeGen2(),
		"azurerm_data_share_dataset_kusto_cluster":  resourceDataShareDataSetKustoCluster(),
		"azurerm_data_share_dataset_kusto_database": resourceDataShareDataSetKustoDatabase(),
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
