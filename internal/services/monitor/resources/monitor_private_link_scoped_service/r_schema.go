// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scoped_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorPrivateLinkScopedServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Azure Monitor Private Link Scoped Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameOptional(),

		"scope_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"resource_group_name", "scope_resource_group_name", "scope_resource_id"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scope_subscription_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"scope_name": {
			Description:  "The name of the Azure Monitor Private Link Scope.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"scope_name", "scope_resource_id"},
			ValidateFunc: validate.PrivateLinkScopeName,
		},

		"scope_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"scope_name", "scope_resource_id"},
			ValidateFunc: privatelinkscopesapis.ValidatePrivateLinkScopeID,
		},

		"linked_resource_id": {
			Description: "The ID of the linked resource.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				componentsapis.ValidateComponentID,
				workspaces.ValidateWorkspaceID,
				datacollectionendpoints.ValidateDataCollectionEndpointID,
			),
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}
