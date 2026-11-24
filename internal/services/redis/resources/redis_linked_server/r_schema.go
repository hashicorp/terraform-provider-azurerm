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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"linked_redis_cache_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: linkedserver.ValidateRediID,
		},

		"linked_redis_cache_location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_role": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(linkedserver.PossibleValuesForReplicationRole(), false),
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"geo_replicated_primary_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
