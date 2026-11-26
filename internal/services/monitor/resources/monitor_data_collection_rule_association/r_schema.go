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
			Description:  "The ID of the Azure Resource which to associate to a Data Collection Rule or a Data Collection Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"name": {
			Description: "The name which should be used for this Data Collection Rule Association.",
			// TODO: should this be hard-coded in the Create?
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "configurationAccessEndpoint",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_collection_endpoint_id": {
			Description:  "The ID of the Data Collection Endpoint which will be associated to the target resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: datacollectionendpoints.ValidateDataCollectionEndpointID,
			ExactlyOneOf: []string{"data_collection_endpoint_id", "data_collection_rule_id"},
		},

		"data_collection_rule_id": {
			Description:  "The ID of the Data Collection Rule which will be associated to the target resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: datacollectionrules.ValidateDataCollectionRuleID,
			RequiredWith: []string{"name"},
			ExactlyOneOf: []string{"data_collection_endpoint_id", "data_collection_rule_id"},
		},

		"description": {
			Description: "The description of the Data Collection Rule Association.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
