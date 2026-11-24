// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_linked_service

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogAnalyticsLinkedServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"workspace_id": {
			Description:  "The ID of the Log Analytics Workspace that will contain the Log Analytics Linked Service resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: linkedservices.ValidateWorkspaceID,
		},

		"read_access_id": {
			Description:  "The ID of the readable Resource that will be linked to the workspace. This should be used for linking to an Automation Account resource.",
			Type:         pluginsdk.TypeString,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			ExactlyOneOf: []string{"read_access_id", "write_access_id"},
		},

		"write_access_id": {
			Description:  "The ID of the writable Resource that will be linked to the workspace. This should be used for linking to a Log Analytics Cluster resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			ExactlyOneOf: []string{"read_access_id", "write_access_id"},
		},
		// Exported properties
		"name": {
			Description: "The generated name of the Linked Service. The format for this attribute is always `<workspace name>/<linked service type>`(e.g. `workspace1/Automation` or `workspace1/Cluster`).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
