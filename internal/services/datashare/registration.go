// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/resources/data_share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/resources/data_share_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/resources/data_share_dataset_blob_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/resources/data_share_dataset_data_lake_gen2"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/resources/data_share_dataset_kusto_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/resources/data_share_dataset_kusto_database"
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
		"azurerm_data_share":                        data_share.RegisterDataSource(),
		"azurerm_data_share_account":                data_share_account.RegisterDataSource(),
		"azurerm_data_share_dataset_blob_storage":   data_share_dataset_blob_storage.RegisterDataSource(),
		"azurerm_data_share_dataset_data_lake_gen2": data_share_dataset_data_lake_gen2.RegisterDataSource(),
		"azurerm_data_share_dataset_kusto_cluster":  data_share_dataset_kusto_cluster.RegisterDataSource(),
		"azurerm_data_share_dataset_kusto_database": data_share_dataset_kusto_database.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_share":                        data_share.RegisterResource(),
		"azurerm_data_share_account":                data_share_account.RegisterResource(),
		"azurerm_data_share_dataset_blob_storage":   data_share_dataset_blob_storage.RegisterResource(),
		"azurerm_data_share_dataset_data_lake_gen2": data_share_dataset_data_lake_gen2.RegisterResource(),
		"azurerm_data_share_dataset_kusto_cluster":  data_share_dataset_kusto_cluster.RegisterResource(),
		"azurerm_data_share_dataset_kusto_database": data_share_dataset_kusto_database.RegisterResource(),
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
