// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_linked_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/linkedserver"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func redisLinkedServerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"target_redis_cache_name": {
			Description:  "The name of the Redis Cache to link with (e.g. the primary role).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"linked_redis_cache_id": {
			Description:  "The ID of the linked Redis Cache.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: linkedserver.ValidateRediID,
		},

		"linked_redis_cache_location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_role": {
			Description:  "The role of the linked Redis Cache (e.g. `Secondary`).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(linkedserver.PossibleValuesForReplicationRole(), false),
		},

		"name": {
			Description: "The name of the linked server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"geo_replicated_primary_host_name": {
			Description: "The geo-replicated primary hostname for this linked server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
