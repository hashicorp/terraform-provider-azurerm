// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package applicationgateway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
)

// Use the production listener schema and hash while exercising the planning
// RPC without configuring an Azure client or provisioning a gateway.
func testProvider() (*schema.Provider, *schema.Resource) {
	gateway := (network.Registration{}).SupportedResources()[resourceName]
	resource := &schema.Resource{Schema: map[string]*schema.Schema{"http_listener": gateway.Schema["http_listener"]}}
	provider := &schema.Provider{ResourcesMap: map[string]*schema.Resource{resourceName: resource, "azurerm_other": resource}}
	return provider, resource
}
