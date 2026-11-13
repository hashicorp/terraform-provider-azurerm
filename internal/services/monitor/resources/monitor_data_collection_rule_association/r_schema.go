// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"name": {
			// TODO: should this be hard-coded in the Create?
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "configurationAccessEndpoint",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_collection_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: datacollectionendpoints.ValidateDataCollectionEndpointID,
			ExactlyOneOf: []string{"data_collection_endpoint_id", "data_collection_rule_id"},
		},

		"data_collection_rule_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: datacollectionrules.ValidateDataCollectionRuleID,
			ExactlyOneOf: []string{"data_collection_endpoint_id", "data_collection_rule_id"},
			RequiredWith: []string{"name"},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
