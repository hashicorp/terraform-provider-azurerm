// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Geo-replication / linked databases are managed as a separate resource because when dbs are linked, ARM
// will mutate the state of all, causing unexpected plan diff

type ManagedRedisDatabaseGeoReplicationResource struct{}

var (
	_ sdk.ResourceWithCustomizeDiff = ManagedRedisDatabaseGeoReplicationResource{}
	_ sdk.ResourceWithUpdate        = ManagedRedisDatabaseGeoReplicationResource{}
)

type ManagedRedisDatabaseGeoReplicationResourceModel struct {
	DatabaseId        string   `tfschema:"database_id"`
	LinkedDatabaseIds []string `tfschema:"linked_database_ids"`
}

func (r ManagedRedisDatabaseGeoReplicationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: databases.ValidateDatabaseID,
		},

		"linked_database_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 5,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: databases.ValidateDatabaseID,
			},
		},
	}
}

func (r ManagedRedisDatabaseGeoReplicationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedRedisDatabaseGeoReplicationResource) ModelObject() interface{} {
	return &ManagedRedisDatabaseGeoReplicationResourceModel{}
}

func (r ManagedRedisDatabaseGeoReplicationResource) ResourceType() string {
	return "azurerm_managed_redis_database_geo_replication"
}

func (r ManagedRedisDatabaseGeoReplicationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return databases.ValidateDatabaseID
}

func (r ManagedRedisDatabaseGeoReplicationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			var model ManagedRedisDatabaseGeoReplicationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := databases.ParseDatabaseID(model.DatabaseId)
			if err != nil {
				return err
			}

			if err := linkUnlinkGeoReplication(ctx, client, model, id); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedRedisDatabaseGeoReplicationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			id, err := databases.ParseDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ManagedRedisDatabaseGeoReplicationResourceModel{
				DatabaseId: id.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil && props.GeoReplication != nil {
					state.LinkedDatabaseIds = make([]string, 0, len(pointer.From(props.GeoReplication.LinkedDatabases)))
					for _, db := range pointer.From(props.GeoReplication.LinkedDatabases) {
						state.LinkedDatabaseIds = append(state.LinkedDatabaseIds, pointer.From(db.Id))
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedRedisDatabaseGeoReplicationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			id, err := databases.ParseDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedRedisDatabaseGeoReplicationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := linkUnlinkGeoReplication(ctx, client, model, id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r ManagedRedisDatabaseGeoReplicationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// Deletion works by unlinking all databases except self

			client := metadata.Client.ManagedRedis.DatabaseClient

			id, err := databases.ParseDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			if existing.Model.Properties != nil && existing.Model.Properties.GeoReplication != nil {
				fromDbIds := flattenLinkedDatabases(existing.Model.Properties.GeoReplication.LinkedDatabases)
				toDbIds := []string{id.ID()}

				dbIdsToUnlink := databaselink.DbIdsToUnlink(fromDbIds, toDbIds)

				if len(dbIdsToUnlink) > 0 {
					params := databases.ForceUnlinkParameters{
						Ids: dbIdsToUnlink,
					}
					if err := client.ForceUnlinkThenPoll(ctx, *id, params); err != nil {
						return fmt.Errorf("force unlink %s: %+v", *id, err)
					}
				}
			}

			return nil
		},
	}
}

func (r ManagedRedisDatabaseGeoReplicationResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var model ManagedRedisDatabaseGeoReplicationResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return err
			}

			if model.DatabaseId != "" && len(model.LinkedDatabaseIds) > 0 {
				if linkedDbsIncludeSelf := slices.Contains(model.LinkedDatabaseIds, model.DatabaseId); !linkedDbsIncludeSelf {
					return fmt.Errorf("`linked_database_ids` must include `database_id`: %s", model.DatabaseId)
				}
			}

			return nil
		},
	}
}

func linkUnlinkGeoReplication(ctx context.Context, client *databases.DatabasesClient, model ManagedRedisDatabaseGeoReplicationResourceModel, id *databases.DatabaseId) error {
	existing, err := client.Get(ctx, *id)
	if err != nil {
		return err
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}
	if existing.Model.Properties.GeoReplication == nil {
		return fmt.Errorf("geo_replication_group_name has to be set on database %s", *id)
	}

	fromDbIds := flattenLinkedDatabases(existing.Model.Properties.GeoReplication.LinkedDatabases)
	toDbIds := model.LinkedDatabaseIds

	// Additional validation to improve error message, otherwise user will get generic 400 error message.
	// This validation often gets skipped in CustomizeDiff because the value is not known before apply
	if linkedDbsIncludeSelf := slices.Contains(toDbIds, model.DatabaseId); !linkedDbsIncludeSelf {
		return fmt.Errorf("`linked_database_ids` must include `database_id`: %s", model.DatabaseId)
	}

	dbIdsToUnlink := databaselink.DbIdsToUnlink(fromDbIds, toDbIds)

	if len(dbIdsToUnlink) > 0 {
		params := databases.ForceUnlinkParameters{
			Ids: dbIdsToUnlink,
		}
		if err := client.ForceUnlinkThenPoll(ctx, *id, params); err != nil {
			return fmt.Errorf("force unlink %s: %+v", *id, err)
		}
	}

	if databaselink.HasDbToLink(fromDbIds, toDbIds) {
		params := databases.ForceLinkParameters{
			GeoReplication: databases.ForceLinkParametersGeoReplication{
				GroupNickname:   existing.Model.Properties.GeoReplication.GroupNickname,
				LinkedDatabases: expandLinkedDatabases(toDbIds),
			},
		}

		err = client.ForceLinkToReplicationGroupThenPoll(ctx, *id, params)
		if err != nil {
			return fmt.Errorf("force link %s: %+v", *id, err)
		}
	}
	return nil
}

func flattenLinkedDatabases(dbs *[]databases.LinkedDatabase) []string {
	if dbs == nil {
		return nil
	}

	result := make([]string, 0, len(*dbs))
	for _, db := range *dbs {
		if db.Id != nil {
			result = append(result, pointer.From(db.Id))
		}
	}
	return result
}

func expandLinkedDatabases(dbIds []string) *[]databases.LinkedDatabase {
	if len(dbIds) == 0 {
		return nil
	}

	result := make([]databases.LinkedDatabase, 0, len(dbIds))
	for _, id := range dbIds {
		result = append(result, databases.LinkedDatabase{
			Id: pointer.To(id),
		})
	}
	return &result
}
