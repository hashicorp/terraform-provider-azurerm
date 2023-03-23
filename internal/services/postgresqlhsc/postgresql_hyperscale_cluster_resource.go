package postgresqlhsc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgresqlhsc/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLHyperScaleClusterModel struct {
	Name                             string              `tfschema:"name"`
	ResourceGroupName                string              `tfschema:"resource_group_name"`
	Location                         string              `tfschema:"location"`
	AdministratorLoginPassword       string              `tfschema:"administrator_login_password"`
	CitusVersion                     string              `tfschema:"citus_version"`
	CoordinatorPublicIPAccessEnabled bool                `tfschema:"coordinator_public_ip_access_enabled"`
	CoordinatorServerEdition         string              `tfschema:"coordinator_server_edition"`
	CoordinatorStorageQuotaInMb      int64               `tfschema:"coordinator_storage_quota_in_mb"`
	CoordinatorVCores                int64               `tfschema:"coordinator_vcores"`
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

type PostgreSQLHyperScaleClusterResource struct{}

var _ sdk.ResourceWithUpdate = PostgreSQLHyperScaleClusterResource{}

func (r PostgreSQLHyperScaleClusterResource) ResourceType() string {
	return "azurerm_postgresql_hyperscale_cluster"
}

func (r PostgreSQLHyperScaleClusterResource) ModelObject() interface{} {
	return &PostgreSQLHyperScaleClusterModel{}
}

func (r PostgreSQLHyperScaleClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return clusters.ValidateServerGroupsv2ID
}

func (r PostgreSQLHyperScaleClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validate.AdministratorLoginPassword,
		},

		"coordinator_storage_quota_in_mb": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ValidateFunc: validation.All(
				validation.IntBetween(32768, 16777216),
				validation.IntDivisibleBy(1024),
			),
		},

		"coordinator_vcores": {
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
			}, false),
		},

		"coordinator_public_ip_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"coordinator_server_edition": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "GeneralPurpose",
			ValidateFunc: validation.StringIsNotEmpty,
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

		"node_server_edition": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "MemoryOptimized",
			ValidateFunc: validation.StringIsNotEmpty,
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
			RequiredWith:     []string{"source_resource_id", "point_in_time_in_utc"},
		},

		"source_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateServerGroupsv2ID,
			RequiredWith: []string{"source_location", "point_in_time_in_utc"},
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

func (r PostgreSQLHyperScaleClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"earliest_restore_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r PostgreSQLHyperScaleClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PostgreSQLHyperScaleClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PostgreSQLHSC.ClustersClient
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
					CoordinatorVCores:               &model.CoordinatorVCores,
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
			if v := metadata.ResourceData.GetRawConfig().AsValueMap()["shards_on_coordinator_enabled"]; !v.IsNull() {
				parameters.Properties.EnableShardsOnCoordinator = &model.ShardsOnCoordinatorEnabled
			}

			if v := model.Tags; v != nil {
				parameters.Tags = &v
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PostgreSQLHyperScaleClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.ClustersClient

			id, err := clusters.ParseServerGroupsv2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PostgreSQLHyperScaleClusterModel
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

			if metadata.ResourceData.HasChange("coordinator_vcores") {
				parameters.Properties.CoordinatorVCores = &model.CoordinatorVCores
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

func (r PostgreSQLHyperScaleClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.ClustersClient

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

			state := PostgreSQLHyperScaleClusterModel{
				Name:              id.ServerGroupsv2Name,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if props := model.Properties; props != nil {
				state.AdministratorLoginPassword = metadata.ResourceData.Get("administrator_login_password").(string)
				state.SourceResourceId = metadata.ResourceData.Get("source_resource_id").(string)
				state.SourceLocation = metadata.ResourceData.Get("source_location").(string)
				state.PointInTimeInUTC = metadata.ResourceData.Get("point_in_time_in_utc").(string)
				state.CoordinatorPublicIPAccessEnabled = *props.CoordinatorEnablePublicIPAccess
				state.CoordinatorServerEdition = *props.CoordinatorServerEdition
				state.CoordinatorStorageQuotaInMb = *props.CoordinatorStorageQuotaInMb
				state.CoordinatorVCores = *props.CoordinatorVCores
				state.HaEnabled = *props.EnableHa
				state.NodeCount = *props.NodeCount
				state.NodePublicIPAccessEnabled = *props.NodeEnablePublicIPAccess
				state.NodeServerEdition = *props.NodeServerEdition
				state.NodeStorageQuotaInMb = *props.NodeStorageQuotaInMb
				state.NodeVCores = *props.NodeVCores
				state.ShardsOnCoordinatorEnabled = *props.EnableShardsOnCoordinator

				if v := props.CitusVersion; v != nil {
					state.CitusVersion = *v
				}

				if v := props.MaintenanceWindow; v != nil {
					state.MaintenanceWindow = flattenMaintenanceWindow(v)
				}

				if v := props.PreferredPrimaryZone; v != nil {
					state.PreferredPrimaryZone = *v
				}

				if v := props.PostgresqlVersion; v != nil {
					state.SqlVersion = *v
				}

				if v := props.EarliestRestoreTime; v != nil {
					state.EarliestRestoreTime = *v
				}
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PostgreSQLHyperScaleClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.ClustersClient

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
			CustomWindow: utils.String("Disabled"),
		}
	}

	v := input[0]

	maintenanceWindow := clusters.MaintenanceWindow{
		CustomWindow: utils.String("Enabled"),
		StartHour:    utils.Int64(v.StartHour),
		StartMinute:  utils.Int64(v.StartMinute),
		DayOfWeek:    utils.Int64(v.DayOfWeek),
	}

	return &maintenanceWindow
}

func flattenMaintenanceWindow(input *clusters.MaintenanceWindow) []MaintenanceWindow {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == "Disabled" {
		return nil
	}

	result := MaintenanceWindow{}

	if input.DayOfWeek != nil {
		result.DayOfWeek = *input.DayOfWeek
	}

	if input.StartHour != nil {
		result.StartHour = *input.StartHour
	}

	if input.StartMinute != nil {
		result.StartMinute = *input.StartMinute
	}

	return []MaintenanceWindow{result}
}
