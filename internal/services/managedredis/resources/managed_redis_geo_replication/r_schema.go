// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_redis_geo_replication

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_redis_id": {
			Description:  "The ID of the Managed Redis through which geo-replication group will be managed. Linking is reciprocal, if A is linked to B, both A and B will have the same linking state. There is no need to have duplicate `azurerm_managed_redis_geo_replication` resources for each.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
		},

		"linked_managed_redis_ids": {
			Description: "A set of other Managed Redis IDs to link together in the geo-replication group. The ID of this Managed Redis is always included by default and does not need to be provided here. Can contain up to 4 Managed Redis IDs, making up a group of 5 in total. All Managed Redis must have the same `geo_replication_group_name` configured. Once linked, the geo-replication state of all Managed Redis will be updated.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			MaxItems:    4,
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

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var model ManagedRedisGeoReplicationResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return err
			}

			if model.ManagedRedisId != "" && slices.ContainsFunc(model.LinkedManagedRedisIds, func(id string) bool {
				return id != "" && id == model.ManagedRedisId
			}) {
				return fmt.Errorf("linked_managed_redis_ids cannot contain the same value as managed_redis_id: %s", model.ManagedRedisId)
			}

			return nil
		},
	}
}
