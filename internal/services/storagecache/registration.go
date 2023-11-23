// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}
var _ sdk.UntypedServiceRegistration = Registration{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Storage Cache"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Azure Managed Lustre File System",
		"Storage",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_hpc_cache":                 resourceHPCCache(),
		"azurerm_hpc_cache_access_policy":   resourceHPCCacheAccessPolicy(),
		"azurerm_hpc_cache_blob_target":     resourceHPCCacheBlobTarget(),
		"azurerm_hpc_cache_blob_nfs_target": resourceHPCCacheBlobNFSTarget(),
		"azurerm_hpc_cache_nfs_target":      resourceHPCCacheNFSTarget(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ManagedLustreFileSystemResource{},
	}
}
