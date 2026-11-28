// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Azure Relay Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(6, 50),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Description:  "The name of the SKU to use.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForSkuName(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"metric_id": {
			Description: "The Identifier for Azure Insights metrics.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"primary_connection_string": {
			Description: "The primary connection string for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string": {
			Description: "The secondary connection string for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_key": {
			Description: "The primary access key for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_key": {
			Description: "The secondary access key for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
