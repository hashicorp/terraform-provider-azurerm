// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_replication

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			var model ResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := redisenterprise.ParseRedisEnterpriseID(model.ManagedRedisId)
			if err != nil {
				return err
			}

			if err := linkUnlinkGeoReplication(ctx, metadata, client, model, id, true); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var model ResourceModel
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
