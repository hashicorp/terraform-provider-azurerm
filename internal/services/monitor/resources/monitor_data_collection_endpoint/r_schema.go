// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				datacollectionendpoints.PossibleValuesForKnownDataCollectionEndpointResourceKind(), false,
			),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_access_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"immutable_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"logs_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metrics_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
