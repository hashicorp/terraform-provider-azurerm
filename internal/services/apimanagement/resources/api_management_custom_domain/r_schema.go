// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_custom_domain

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementCustomDomainSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"management": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
			Elem: &pluginsdk.Resource{
				Schema: helpers.ApiManagementResourceHostnameSchema(),
			},
		},
		"portal": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
			Elem: &pluginsdk.Resource{
				Schema: helpers.ApiManagementResourceHostnameSchema(),
			},
		},
		"developer_portal": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
			Elem: &pluginsdk.Resource{
				Schema: helpers.ApiManagementResourceHostnameSchema(),
			},
		},
		"scm": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
			Elem: &pluginsdk.Resource{
				Schema: helpers.ApiManagementResourceHostnameSchema(),
			},
		},
		"gateway": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"management", "portal", "developer_portal", "gateway", "scm"},
			Elem: &pluginsdk.Resource{
				Schema: helpers.ApiManagementResourceHostnameProxySchema(),
			},
		},
	}
}
