// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub_namespace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func notificationHubNamespaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name to use for this Notification Hub Namespace.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Description:  "The name of the SKU to use for this Notification Hub Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForSkuName(), false),
		},

		"enabled": {
			Description: "Is this Notification Hub Namespace enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},

		"namespace_type": {
			Description:  "The Type of Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForNamespaceType(), false),
		},

		"zone_redundancy_enabled": {
			Description: "Is Zone Redundancy Enabled for the Notification Hub Namespace?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},

		"replication_region": {
			Description:      "The allowed Replication Region for the Notification Hub Namespace.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			Default:          namespaces.ReplicationRegionDefault,
			ValidateFunc:     validation.StringInSlice(namespaces.PossibleValuesForReplicationRegion(), true),
			DiffSuppressFunc: location.DiffSuppressFunc,
		},

		"tags": commonschema.Tags(),

		"servicebus_endpoint": {
			Description: "The ServiceBus Endpoint for this Notification Hub Namespace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
