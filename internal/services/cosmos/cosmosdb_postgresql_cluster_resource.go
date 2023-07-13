// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var CosmosDbPostgreSQLClusterResourceName = "azurerm_cosmosdb_postgresql_cluster"

type CosmosDbPostgreSQLClusterModel struct {
	Name                             string              `tfschema:"name"`
	ResourceGroupName                string              `tfschema:"resource_group_name"`
	Location                         string              `tfschema:"location"`
	AdministratorLoginPassword       string              `tfschema:"administrator_login_password"`
	CitusVersion                     string              `tfschema:"citus_version"`
	CoordinatorPublicIPAccessEnabled bool                `tfschema:"coordinator_public_ip_access_enabled"`
	CoordinatorServerEdition         string              `tfschema:"coordinator_server_edition"`
	CoordinatorStorageQuotaInMb      int64               `tfschema:"coordinator_storage_quota_in_mb"`
	CoordinatorVCoreCount            int64               `tfschema:"coordinator_vcore_count"`
	HaEnabled                        bool                `tfschema:"ha_enabled"`
	ShardsOnCoordinatorEnabled       bool                `tfschema:"shards_on_coordinator_enabled"`
	SourceLocation                   string              `tfschema:"source_location"`
	SourceResourceId                 string              `tfschema:"source_resource_id"`
	MaintenanceWindow                []MaintenanceWindow `tfschema:"maintenance_window"`
	NodeCount                        int64               `tfschema:"node_count"`
	NodePublicIPAccessEnabled        bool                `tfschema:"node_public_ip_access_enabled"`
	NodeServerEdition                string              `tfschema:"node_server_edition"`
	NodeStorageQuotaInMb             int64               `tfschema:"node_storage_quota_in_mb"`
	NodeVCores                       int64               `tfschema:"node_vcores"`
	PointInTimeInUTC                 string              `tfschema:"point_in_time_in_utc"`
	PreferredPrimaryZone             string              `tfschema:"preferred_primary_zone"`
	SqlVersion                       string              `tfschema:"sql_version"`
	Tags                             map[string]string   `tfschema:"tags"`
	EarliestRestoreTime              string              `tfschema:"earliest_restore_time"`
}

type MaintenanceWindow struct {
	DayOfWeek   int64 `tfschema:"day_of_week"`
	StartHour   int64 `tfschema:"start_hour"`
	StartMinute int64 `tfschema:"start_minute"`
}

type CosmosDbPostgreSQLClusterResource struct{}

var _ sdk.ResourceWithUpdate = CosmosDbPostgreSQLClusterResource{}

func (r CosmosDbPostgreSQLClusterResource) ResourceType() string {
	return CosmosDbPostgreSQLClusterResourceName
}

func (r CosmosDbPostgreSQLClusterResource) ModelObject() interface{} {
	return &CosmosDbPostgreSQLClusterModel{}
}

func (r CosmosDbPostgreSQLClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return clusters.ValidateServerGroupsv2ID
}

func (r CosmosDbPostgreSQLClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 260),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(8, 256),
		},

		"coordinator_storage_quota_in_mb": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ValidateFunc: validation.All(
				validation.IntBetween(32768, 16777216),
				validation.IntDivisibleBy(1024),
			),
		},

		"coordinator_vcore_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ValidateFunc: validation.IntInSlice([]int{
				1,
				2,
				4,
				8,
				16,
				32,
				64,
				96,
			}),
		},

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ValidateFunc: validation.All(
				validation.IntBetween(0, 20),
				validation.IntNotInSlice([]int{1}),
			),
		},

		"citus_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				"8.3",
				"9.0",
				"9.1",
				"9.2",
				"9.3",
				"9.4",
				"9.5",
				"10.0",
				"10.1",
				"10.2",
				"11.0",
				"11.1",
				"11.2",
				"11.3",
			}, false),
		},

		"coordinator_public_ip_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		// Once the issue https://github.com/Azure/azure-rest-api-specs/issues/23317 is fixed, we would submit PR to improve the validation
		"coordinator_server_edition": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "GeneralPurpose",
			ValidateFunc: validation.StringInSlice([]string{
				"BurstableGeneralPurpose",
				"BurstableMemoryOptimized",
				"GeneralPurpose",
				"MemoryOptimized",
			}, false),
		},

		"ha_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day_of_week": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 6),
					},

					"start_hour": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 23),
					},

					"start_minute": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"node_public_ip_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// Once the issue https://github.com/Azure/azure-rest-api-specs/issues/23317 is fixed, we would submit PR to improve the validation
		"node_server_edition": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "MemoryOptimized",
			ValidateFunc: validation.StringInSlice([]string{
				"BurstableGeneralPurpose",
				"BurstableMemoryOptimized",
				"GeneralPurpose",
				"MemoryOptimized",
			}, false),
		},

		"node_storage_quota_in_mb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.All(
				validation.IntBetween(32768, 16777216),
				validation.IntDivisibleBy(1024),
			),
		},

		"node_vcores": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.IntInSlice([]int{
				1,
				2,
				4,
				8,
				16,
				32,
				64,
				96,
				104,
			}),
		},

		"point_in_time_in_utc": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
			RequiredWith: []string{"source_location", "source_resource_id"},
		},

		"preferred_primary_zone": commonschema.ZoneSingleOptional(),

		"shards_on_coordinator_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"source_location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			RequiredWith:     []string{"source_resource_id"},
		},

		"source_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateServerGroupsv2ID,
			RequiredWith: []string{"source_location"},
		},

		"sql_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				"11",
				"12",
				"13",
				"14",
				"15",
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r CosmosDbPostgreSQLClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"earliest_restore_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CosmosDbPostgreSQLClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbPostgreSQLClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.ClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := clusters.NewServerGroupsv2ID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &clusters.Cluster{
				Location: location.Normalize(model.Location),
				Properties: &clusters.ClusterProperties{
					AdministratorLoginPassword:      &model.AdministratorLoginPassword,
					CoordinatorEnablePublicIPAccess: &model.CoordinatorPublicIPAccessEnabled,
					CoordinatorServerEdition:        &model.CoordinatorServerEdition,
					CoordinatorStorageQuotaInMb:     &model.CoordinatorStorageQuotaInMb,
					CoordinatorVCores:               &model.CoordinatorVCoreCount,
					EnableHa:                        &model.HaEnabled,
					NodeCount:                       &model.NodeCount,
					NodeEnablePublicIPAccess:        &model.NodePublicIPAccessEnabled,
					NodeServerEdition:               &model.NodeServerEdition,
				},
			}

			if v := model.CitusVersion; v != "" {
				parameters.Properties.CitusVersion = &model.CitusVersion
			}

			if v := model.MaintenanceWindow; v != nil {
				parameters.Properties.MaintenanceWindow = expandMaintenanceWindow(v)
			}

			if v := model.NodeStorageQuotaInMb; v != 0 {
				parameters.Properties.NodeStorageQuotaInMb = utils.Int64(model.NodeStorageQuotaInMb)
			}

			if v := model.NodeVCores; v != 0 {
				parameters.Properties.NodeVCores = utils.Int64(model.NodeVCores)
			}

			if v := model.PointInTimeInUTC; v != "" {
				parameters.Properties.PointInTimeUTC = &model.PointInTimeInUTC
			}

			if v := model.SqlVersion; v != "" {
				parameters.Properties.PostgresqlVersion = &model.SqlVersion
			}

			if v := model.PreferredPrimaryZone; v != "" {
				parameters.Properties.PreferredPrimaryZone = &model.PreferredPrimaryZone
			}

			if v := model.SourceLocation; v != "" {
				parameters.Properties.SourceLocation = &model.SourceLocation
			}

			if v := model.SourceResourceId; v != "" {
				parameters.Properties.SourceResourceId = &model.SourceResourceId
			}

			// If `shards_on_coordinator_enabled` isn't set, API would set it to `true` when `node_count` is `0`.
			// If `shards_on_coordinator_enabled` isn't set, API would set it to `false` when `node_count` is greater than or equal to `2`.
			// As `shards_on_coordinator_enabled` is `bool` and it's always set to `false` as zero value when it isn't set, so we cannot use `model.ShardsOnCoordinatorEnabled` to check if this property is set in tf config.
			// nolint staticcheck
			if v, ok := metadata.ResourceData.GetOkExists("shards_on_coordinator_enabled"); ok {
				parameters.Properties.EnableShardsOnCoordinator = utils.Bool(v.(bool))
			}

			if v := model.Tags; v != nil {
				parameters.Tags = &v
			}

			if err := client.CreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CosmosDbPostgreSQLClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ClustersClient

			id, err := clusters.ParseServerGroupsv2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CosmosDbPostgreSQLClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := clusters.ClusterForUpdate{
				Properties: &clusters.ClusterPropertiesForUpdate{},
			}

			if metadata.ResourceData.HasChange("administrator_login_password") {
				parameters.Properties.AdministratorLoginPassword = &model.AdministratorLoginPassword
			}

			if metadata.ResourceData.HasChange("citus_version") {
				parameters.Properties.CitusVersion = &model.CitusVersion
			}

			if metadata.ResourceData.HasChange("coordinator_public_ip_access_enabled") {
				parameters.Properties.CoordinatorEnablePublicIPAccess = &model.CoordinatorPublicIPAccessEnabled
			}

			if metadata.ResourceData.HasChange("coordinator_server_edition") {
				parameters.Properties.CoordinatorServerEdition = &model.CoordinatorServerEdition
			}

			if metadata.ResourceData.HasChange("coordinator_storage_quota_in_mb") {
				parameters.Properties.CoordinatorStorageQuotaInMb = &model.CoordinatorStorageQuotaInMb
			}

			if metadata.ResourceData.HasChange("coordinator_vcore_count") {
				parameters.Properties.CoordinatorVCores = &model.CoordinatorVCoreCount
			}

			if metadata.ResourceData.HasChange("ha_enabled") {
				parameters.Properties.EnableHa = &model.HaEnabled
			}

			if metadata.ResourceData.HasChange("maintenance_window") {
				parameters.Properties.MaintenanceWindow = expandMaintenanceWindow(model.MaintenanceWindow)
			}

			if metadata.ResourceData.HasChange("node_count") {
				parameters.Properties.NodeCount = &model.NodeCount
			}

			if metadata.ResourceData.HasChange("node_public_ip_access_enabled") {
				parameters.Properties.NodeEnablePublicIPAccess = &model.NodePublicIPAccessEnabled
			}

			if metadata.ResourceData.HasChange("node_server_edition") {
				parameters.Properties.NodeServerEdition = &model.NodeServerEdition
			}

			if metadata.ResourceData.HasChange("node_storage_quota_in_mb") {
				parameters.Properties.NodeStorageQuotaInMb = utils.Int64(model.NodeStorageQuotaInMb)
			}

			if metadata.ResourceData.HasChange("node_vcores") {
				parameters.Properties.NodeVCores = utils.Int64(model.NodeVCores)
			}

			if metadata.ResourceData.HasChange("preferred_primary_zone") {
				parameters.Properties.PreferredPrimaryZone = &model.PreferredPrimaryZone
			}

			if metadata.ResourceData.HasChange("shards_on_coordinator_enabled") {
				parameters.Properties.EnableShardsOnCoordinator = &model.ShardsOnCoordinatorEnabled
			}

			if metadata.ResourceData.HasChange("sql_version") {
				parameters.Properties.PostgresqlVersion = &model.SqlVersion
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CosmosDbPostgreSQLClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ClustersClient

			id, err := clusters.ParseServerGroupsv2ID(metadata.ResourceData.Id())
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CosmosDbPostgreSQLClusterModel{
				Name:              id.ServerGroupsv2Name,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if props := model.Properties; props != nil {
				state.AdministratorLoginPassword = metadata.ResourceData.Get("administrator_login_password").(string)
				state.SourceResourceId = metadata.ResourceData.Get("source_resource_id").(string)
				state.SourceLocation = metadata.ResourceData.Get("source_location").(string)
				state.PointInTimeInUTC = metadata.ResourceData.Get("point_in_time_in_utc").(string)
				state.CoordinatorPublicIPAccessEnabled = pointer.From(props.CoordinatorEnablePublicIPAccess)
				state.CoordinatorServerEdition = pointer.From(props.CoordinatorServerEdition)
				state.CoordinatorStorageQuotaInMb = pointer.From(props.CoordinatorStorageQuotaInMb)
				state.CoordinatorVCoreCount = pointer.From(props.CoordinatorVCores)
				state.HaEnabled = pointer.From(props.EnableHa)
				state.NodeCount = pointer.From(props.NodeCount)
				state.NodePublicIPAccessEnabled = pointer.From(props.NodeEnablePublicIPAccess)
				state.NodeServerEdition = pointer.From(props.NodeServerEdition)
				state.NodeStorageQuotaInMb = pointer.From(props.NodeStorageQuotaInMb)
				state.NodeVCores = pointer.From(props.NodeVCores)
				state.ShardsOnCoordinatorEnabled = pointer.From(props.EnableShardsOnCoordinator)
				state.CitusVersion = pointer.From(props.CitusVersion)
				state.PreferredPrimaryZone = pointer.From(props.PreferredPrimaryZone)
				state.SqlVersion = pointer.From(props.PostgresqlVersion)
				state.EarliestRestoreTime = pointer.From(props.EarliestRestoreTime)

				if v := props.MaintenanceWindow; v != nil {
					state.MaintenanceWindow = flattenMaintenanceWindow(v)
				}
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CosmosDbPostgreSQLClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ClustersClient

			id, err := clusters.ParseServerGroupsv2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandMaintenanceWindow(input []MaintenanceWindow) *clusters.MaintenanceWindow {
	if len(input) == 0 {
		return &clusters.MaintenanceWindow{
			CustomWindow: pointer.To("Disabled"),
		}
	}

	v := input[0]

	maintenanceWindow := clusters.MaintenanceWindow{
		CustomWindow: pointer.To("Enabled"),
		StartHour:    pointer.To(v.StartHour),
		StartMinute:  pointer.To(v.StartMinute),
		DayOfWeek:    pointer.To(v.DayOfWeek),
	}

	return &maintenanceWindow
}

func flattenMaintenanceWindow(input *clusters.MaintenanceWindow) []MaintenanceWindow {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == "Disabled" {
		return nil
	}

	return []MaintenanceWindow{
		{
			DayOfWeek:   pointer.From(input.DayOfWeek),
			StartHour:   pointer.From(input.StartHour),
			StartMinute: pointer.From(input.StartMinute),
		},
	}
}
