// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_replication

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_redis_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
		},

		"linked_managed_redis_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 4,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
