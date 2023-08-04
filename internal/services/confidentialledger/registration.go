// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confidentialledger

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Registration type for Azure Confidential Ledger.
type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Confidential Ledger"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Confidential Ledger",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_confidential_ledger": dataSourceConfidentialLedger(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_confidential_ledger": resourceConfidentialLedger(),
	}
}
