// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/passwordauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/pools"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -test-expect-non-empty
type HorizondbClusterResource struct{}

type HorizondbClusterResourceModel struct {
	Name                           string                               `tfschema:"name"`
	ResourceGroupName              string                               `tfschema:"resource_group_name"`
	Location                       string                               `tfschema:"location"`
	PostgresqlVersion              string                               `tfschema:"postgresql_version"`
	ReadWritePool                  []HorizondbClusterPool               `tfschema:"read_write_pool"`
	AdministratorLogin             string                               `tfschema:"administrator_login"`
	AdministratorPassword          string                               `tfschema:"administrator_password"`
	AdministratorPasswordWOVersion int64                                `tfschema:"administrator_password_wo_version"`
	BackupRetentionInDays          int64                                `tfschema:"backup_retention_in_days"`
	GeoRedundantBackupEnabled      bool                                 `tfschema:"geo_redundant_backup_enabled"`
	Identity                       []identity.ModelUserAssigned         `tfschema:"identity"`
	PublicNetworkAccess            string                               `tfschema:"public_network_access"`
	PointInTimeRestore             []HorizondbClusterPointInTimeRestore `tfschema:"point_in_time_restore"`
	Tags                           map[string]string                    `tfschema:"tags"`
}

type HorizondbClusterPointInTimeRestore struct {
	SourceHorizondbClusterId string `tfschema:"source_horizondb_cluster_id"`
	PointInTimeUtc           string `tfschema:"point_in_time_utc"`
	SourceLocation           string `tfschema:"source_location"`
}

type HorizondbClusterPool struct {
	Name                 string                        `tfschema:"name"`
	Id                   string                        `tfschema:"id"`
	Compute              []HorizondbClusterPoolCompute `tfschema:"compute"`
	NodeCount            int64                         `tfschema:"node_count"`
	HighAvailabilityMode string                        `tfschema:"high_availability_mode"`
}

type HorizondbClusterPoolCompute struct {
	Processor  string `tfschema:"processor"`
	Generation int64  `tfschema:"generation"`
	VCores     int64  `tfschema:"vcores"`
}

var (
	_ sdk.ResourceWithUpdate        = HorizondbClusterResource{}
	_ sdk.ResourceWithIdentity      = HorizondbClusterResource{}
	_ sdk.ResourceWithCustomizeDiff = HorizondbClusterResource{}
)

func (HorizondbClusterResource) ModelObject() any {
	return &HorizondbClusterResourceModel{}
}

func (HorizondbClusterResource) ResourceType() string {
	return "azurerm_horizondb_cluster"
}

func (HorizondbClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return clusters.ValidateClusterID
}

func (HorizondbClusterResource) Identity() resourceids.ResourceId {
	return &clusters.ClusterId{}
}

func (HorizondbClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63),
				validation.StringMatch(regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`), "must contain only lowercase letters, numbers and hyphens and start and end with a letter or number"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 63),
		},

		"read_write_pool": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(1, 63),
							validation.StringMatch(regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`), "must contain only lowercase letters, numbers and hyphens and start and end with a letter or number"),
						),
					},

					"compute": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"processor": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForComputeProcessorType(), false),
								},

								"generation": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntInSlice([]int{1}),
								},

								"vcores": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntInSlice([]int{2, 4, 8, 16, 20, 32, 48, 64, 96}),
								},
							},
						},
					},
					"node_count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 16),
					},

					"high_availability_mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(pools.PossibleValuesForPoolHighAvailabilityMode(), false),
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"administrator_password": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"administrator_password_wo"},
			AtLeastOneOf:  []string{"administrator_password", "administrator_password_wo"},
			ValidateFunc:  validation.StringLenBetween(8, 128),
		},

		"administrator_password_wo": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			ConflictsWith: []string{"administrator_password"},
			AtLeastOneOf:  []string{"administrator_password", "administrator_password_wo"},
			RequiredWith:  []string{"administrator_password_wo_version"},
			ValidateFunc:  validation.StringLenBetween(8, 128),
		},

		"administrator_password_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_password_wo"},
			ValidateFunc: validation.IntAtLeast(1),
		},

		"backup_retention_in_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 35),
		},

		"geo_redundant_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"point_in_time_restore": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"postgresql_version", "point_in_time_restore"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"point_in_time_utc": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},

					"source_horizondb_cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&clusters.ClusterId{}),

					"source_location": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						StateFunc:    location.StateFunc,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"postgresql_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C - Restored clusters inherit their PostgreSQL version from the source backup.
			Computed:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"postgresql_version", "point_in_time_restore"},
			ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForPostgreSqlVersion(), false),
		},

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(clusters.PublicNetworkAccessStateDisabled),
			ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForPublicNetworkAccessState(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (HorizondbClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (HorizondbClusterResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			if err := validateHorizondbClusterPoolConfig(diff.GetRawConfig()); err != nil {
				return err
			}
			return nil
		},
	}
}

func (r HorizondbClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state HorizondbClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := clusters.NewClusterID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.Name)
			// Pool changes also put the cluster into Updating, so serialize mutations using the cluster ID.
			lockId := strings.ToLower(id.ID())
			locks.ByID(lockId)
			defer locks.UnlockByID(lockId)

			client := metadata.Client.HorizonDB.ClustersClient
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			geoRedundancy := clusters.GeoRedundancyStateDisabled
			if state.GeoRedundantBackupEnabled {
				geoRedundancy = clusters.GeoRedundancyStateEnabled
			}

			payload := clusters.Cluster{
				Location: location.Normalize(state.Location),
				Tags:     pointer.To(state.Tags),
				Properties: &clusters.ClusterProperties{
					CreationMode: pointer.To(clusters.ClusterCreateModeCreate),
					Network:      &clusters.Network{PublicAccess: pointer.ToEnum[clusters.PublicNetworkAccessState](state.PublicNetworkAccess)},
					Backup: &clusters.BackupProperties{
						GeoRedundancy: geoRedundancy,
						RetentionDays: pointer.To(state.BackupRetentionInDays),
					},
				},
			}

			if len(state.Identity) > 0 {
				expanded, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expanded
			}

			if len(state.PointInTimeRestore) == 0 {
				payload.Properties.PostgreSqlVersion = pointer.ToEnum[clusters.PostgreSqlVersion](state.PostgresqlVersion)
			} else {
				payload.Properties.CreationMode = pointer.To(clusters.ClusterCreateModePointInTimeRestore)
				restore := state.PointInTimeRestore[0]
				payload.Properties.CreationSourceClusterResourceId = pointer.To(restore.SourceHorizondbClusterId)
				payload.Properties.PointInTimeUtc = pointer.To(restore.PointInTimeUtc)
				if restore.SourceLocation != "" {
					payload.Properties.CreationSourceClusterLocation = pointer.To(location.Normalize(restore.SourceLocation))
				}
			}

			if err := client.CreateOrReplaceCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			if err := createHorizondbClusterPool(ctx, metadata, id, state.ReadWritePool[0]); err != nil {
				return err
			}

			return createHorizondbClusterPasswordAuthentication(ctx, metadata, id, state)
		},
	}
}

func (r HorizondbClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := clusters.ParseClusterIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.HorizonDB.ClustersClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			pool, err := readHorizondbClusterPool(ctx, metadata, *id, resp.Model)
			if err != nil {
				return err
			}

			login, err := readHorizondbClusterAdministratorLogin(ctx, metadata.Client.HorizonDB.PasswordAuthenticationsClient, *id)
			if err != nil {
				return err
			}

			return r.flatten(metadata, id, resp.Model, pool, login)
		},
	}
}

func (HorizondbClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := clusters.ParseClusterIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			lockId := strings.ToLower(id.ID())
			locks.ByID(lockId)
			defer locks.UnlockByID(lockId)

			var state HorizondbClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChanges("identity", "tags", "backup_retention_in_days", "public_network_access") {
				client := metadata.Client.HorizonDB.ClustersClient
				existing, err := client.Get(ctx, *id)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
				if existing.Model == nil {
					return fmt.Errorf("retrieving %s: `model` was nil", *id)
				}

				if existing.Model.Properties == nil {
					return fmt.Errorf("retrieving %s: `properties` was nil", *id)
				}

				props := existing.Model.Properties
				payload := clusters.Cluster{
					Location: existing.Model.Location,
					Identity: existing.Model.Identity,
					Tags:     existing.Model.Tags,
					Properties: &clusters.ClusterProperties{
						PostgreSqlVersion: props.PostgreSqlVersion,
						Backup:            props.Backup,
						Network:           props.Network,
						DataEncryption:    props.DataEncryption,
						AiModelManagement: props.AiModelManagement,
					},
				}

				// Backup retention is in a separate change group and must be updated independently from other property change groups.
				if metadata.ResourceData.HasChange("backup_retention_in_days") {
					if payload.Properties.Backup == nil {
						payload.Properties.Backup = &clusters.BackupProperties{}
					}
					payload.Properties.Backup.RetentionDays = pointer.To(state.BackupRetentionInDays)
					if err := client.CreateOrReplaceThenPoll(ctx, *id, payload); err != nil {
						return fmt.Errorf("updating backup retention for %s: %+v", id, err)
					}
				}

				// Public network access is in a separate change group and must be updated independently from other property change groups.
				if metadata.ResourceData.HasChange("public_network_access") {
					if payload.Properties.Network == nil {
						payload.Properties.Network = &clusters.Network{}
					}
					payload.Properties.Network.PublicAccess = pointer.ToEnum[clusters.PublicNetworkAccessState](state.PublicNetworkAccess)
					if err := client.CreateOrReplaceThenPoll(ctx, *id, payload); err != nil {
						return fmt.Errorf("updating public network access for %s: %+v", id, err)
					}
				}

				// Tags is in a separate change group and must be updated independently from other property change groups.
				if metadata.ResourceData.HasChange("tags") {
					payload.Tags = pointer.To(state.Tags)
				}

				// Aggregate other changes into a single update for the cluster.
				// All other changes that are not part of separate change groups are aggregated here.
				// Service team will fix the mutually exclusive update handling in the future, we can clean this up later.
				if metadata.ResourceData.HasChange("identity") {
					expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
					if err != nil {
						return fmt.Errorf("expanding `identity`: %+v", err)
					}
					payload.Identity = expandedIdentity
				}

				if err := client.CreateOrReplaceThenPoll(ctx, *id, payload); err != nil {
					return fmt.Errorf("updating cluster for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("read_write_pool") {
				if err := updateHorizondbClusterPool(ctx, metadata, *id, state.ReadWritePool[0]); err != nil {
					return err
				}
			}

			if metadata.ResourceData.HasChanges("administrator_password", "administrator_password_wo_version") {
				password := state.AdministratorPassword
				if metadata.ResourceData.HasChange("administrator_password_wo_version") {
					woPassword, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "administrator_password_wo", cty.String)
					if err != nil {
						return err
					}
					if !woPassword.IsNull() {
						password = woPassword.AsString()
					}
				}

				// Changing password can be only sent by PATCH update with only the new password in the payload
				authenticationId := passwordauthentications.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
				payload := passwordauthentications.PasswordAuthenticationForPatchUpdate{
					Properties: &passwordauthentications.PasswordAuthenticationPropertiesForPatchUpdate{
						AdministratorLoginPassword: pointer.To(password),
					},
				}
				if err := metadata.Client.HorizonDB.PasswordAuthenticationsClient.UpdateThenPoll(ctx, authenticationId, payload); err != nil {
					return fmt.Errorf("updating password authentication for %s: %+v", id, err)
				}
			}
			return nil
		},
	}
}

func (HorizondbClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := clusters.ParseClusterIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			lockId := strings.ToLower(id.ID())
			locks.ByID(lockId)
			defer locks.UnlockByID(lockId)

			if err := metadata.Client.HorizonDB.ClustersClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (HorizondbClusterResource) flatten(metadata sdk.ResourceMetaData, id *clusters.ClusterId, model *clusters.Cluster, pool *pools.Pool, administratorLogin string) error {
	if model == nil || model.Properties == nil {
		return nil
	}

	var previous HorizondbClusterResourceModel
	if err := metadata.Decode(&previous); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	props := model.Properties
	state := HorizondbClusterResourceModel{
		Name:                           id.ClusterName,
		ResourceGroupName:              id.ResourceGroupName,
		Location:                       location.Normalize(model.Location),
		PostgresqlVersion:              pointer.FromEnum(props.PostgreSqlVersion),
		AdministratorLogin:             administratorLogin,
		AdministratorPassword:          previous.AdministratorPassword,
		AdministratorPasswordWOVersion: int64(metadata.ResourceData.Get("administrator_password_wo_version").(int)),
		Tags:                           pointer.From(model.Tags),
	}

	if props.Network != nil && props.Network.PublicAccess != nil {
		state.PublicNetworkAccess = pointer.FromEnum(props.Network.PublicAccess)
	}

	if props.Backup != nil {
		state.BackupRetentionInDays = pointer.From(props.Backup.RetentionDays)
		state.GeoRedundantBackupEnabled = props.Backup.GeoRedundancy == clusters.GeoRedundancyStateEnabled
	}

	flattenedIdentity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	state.Identity = pointer.From(flattenedIdentity)

	if pointer.From(props.CreationMode) == clusters.ClusterCreateModePointInTimeRestore {
		restore := HorizondbClusterPointInTimeRestore{
			SourceHorizondbClusterId: pointer.From(props.CreationSourceClusterResourceId),
			SourceLocation:           location.Normalize(pointer.From(props.CreationSourceClusterLocation)),
			PointInTimeUtc:           pointer.From(props.PointInTimeUtc),
		}

		// Azure can omit creation-time restore metadata from GET responses.
		if len(previous.PointInTimeRestore) > 0 {
			if restore.SourceHorizondbClusterId == "" {
				restore.SourceHorizondbClusterId = previous.PointInTimeRestore[0].SourceHorizondbClusterId
			}
			if restore.SourceLocation == "" {
				restore.SourceLocation = location.Normalize(previous.PointInTimeRestore[0].SourceLocation)
			}
			if restore.PointInTimeUtc == "" {
				restore.PointInTimeUtc = previous.PointInTimeRestore[0].PointInTimeUtc
			}
		}

		if restore.SourceHorizondbClusterId != "" || len(previous.PointInTimeRestore) > 0 {
			sourceId, err := clusters.ParseClusterIDInsensitively(restore.SourceHorizondbClusterId)
			if err != nil {
				return err
			}
			restore.SourceHorizondbClusterId = sourceId.ID()
			state.PointInTimeRestore = []HorizondbClusterPointInTimeRestore{restore}
		}
	}

	state.ReadWritePool, err = flattenHorizondbClusterPool(pool)
	if err != nil {
		return err
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func validateHorizondbClusterPoolConfig(config cty.Value) error {
	if !config.IsKnown() || config.IsNull() {
		return nil
	}

	block := config.GetAttr("read_write_pool")
	if !block.IsKnown() || block.IsNull() || block.LengthInt() == 0 {
		return nil
	}

	pool := block.AsValueSlice()[0]
	if !pool.IsKnown() || pool.IsNull() {
		return nil
	}

	name, clusterName := pool.GetAttr("name"), config.GetAttr("name")
	if name.IsKnown() && !name.IsNull() && clusterName.IsKnown() && !clusterName.IsNull() && strings.EqualFold(name.AsString(), clusterName.AsString()) {
		return errors.New("`read_write_pool.name` must differ from the cluster `name`")
	}

	count, mode := pool.GetAttr("node_count"), pool.GetAttr("high_availability_mode")
	if count.IsKnown() && !count.IsNull() && mode.IsKnown() && !mode.IsNull() {
		nodes, _ := count.AsBigFloat().Int64()

		if nodes == 1 && mode.AsString() != string(pools.PoolHighAvailabilityModeDisabled) {
			return errors.New("`read_write_pool.high_availability_mode` must be `Disabled` for one node")
		}

		if nodes > 1 && mode.AsString() == string(pools.PoolHighAvailabilityModeDisabled) {
			return errors.New("`read_write_pool.high_availability_mode` must not be `Disabled` for two or more nodes")
		}
	}

	return nil
}

func createHorizondbClusterPool(ctx context.Context, metadata sdk.ResourceMetaData, clusterId clusters.ClusterId, state HorizondbClusterPool) error {
	id := pools.NewPoolID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ClusterName, state.Name)
	compute := state.Compute[0]
	payload := pools.Pool{
		Properties: &pools.PoolProperties{
			ComputeModel: pools.PoolComputeModel{
				Type:              pools.PoolComputeModelTypeProvisioned,
				ComputeProcessor:  pools.ComputeProcessorType(compute.Processor),
				ComputeGeneration: compute.Generation,
				VCores:            compute.VCores,
			},
			NodeCount:            state.NodeCount,
			HighAvailability:     pools.PoolHighAvailability{Mode: pools.PoolHighAvailabilityMode(state.HighAvailabilityMode)},
			ScheduledMaintenance: pools.SystemScheduledMaintenanceProperties{},
		},
	}

	if err := metadata.Client.HorizonDB.PoolsClient.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating or replacing %s: %+v", id, err)
	}

	return nil
}

func updateHorizondbClusterPool(ctx context.Context, metadata sdk.ResourceMetaData, clusterId clusters.ClusterId, state HorizondbClusterPool) error {
	id := pools.NewPoolID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ClusterName, state.Name)

	client := metadata.Client.HorizonDB.PoolsClient
	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: pool properties are missing", id)
	}

	props := existing.Model.Properties
	payload := pools.Pool{
		Properties: &pools.PoolProperties{
			ComputeModel:     props.ComputeModel,
			HighAvailability: props.HighAvailability,
			NodeCount:        props.NodeCount,
			// ParameterGroup.SyncStatus is read-only and cannot be sent in the update payload.
			ParameterGroup: &pools.PoolParameterGroupReference{
				ApplyMethod:              props.ParameterGroup.ApplyMethod,
				ParameterGroupResourceId: props.ParameterGroup.ParameterGroupResourceId,
			},
			ScheduledMaintenance: props.ScheduledMaintenance,
		},
	}

	// Compute is in separate change group and cannot be updated with HA/node count.
	if metadata.ResourceData.HasChange("read_write_pool.0.compute") {
		compute := state.Compute[0]
		payload.Properties.ComputeModel = pools.PoolComputeModel{
			Type:              pools.PoolComputeModelTypeProvisioned,
			ComputeProcessor:  pools.ComputeProcessorType(compute.Processor),
			ComputeGeneration: compute.Generation,
			VCores:            compute.VCores,
		}
		if err := client.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
			return fmt.Errorf("updating compute for %s: %+v", id, err)
		}
	}

	if metadata.ResourceData.HasChange("read_write_pool.0.high_availability_mode") {
		payload.Properties.HighAvailability.Mode = pools.PoolHighAvailabilityMode(state.HighAvailabilityMode)
	}

	if metadata.ResourceData.HasChange("read_write_pool.0.node_count") {
		payload.Properties.NodeCount = state.NodeCount
	}

	if err := client.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	return nil
}

func readHorizondbClusterPool(ctx context.Context, metadata sdk.ResourceMetaData, id clusters.ClusterId, cluster *clusters.Cluster) (*pools.Pool, error) {
	if cluster == nil || cluster.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: cluster properties are missing", id)
	}

	poolId, err := pools.ParsePoolIDInsensitively(pointer.From(cluster.Properties.ReadWritePool.DesignatedPoolResourceId))
	if err != nil {
		return nil, err
	}

	result, err := metadata.Client.HorizonDB.PoolsClient.Get(ctx, *poolId)
	if err != nil {
		return nil, fmt.Errorf("retrieving designated read/write pool %s: %+v", poolId, err)
	}

	if _, err := flattenHorizondbClusterPool(result.Model); err != nil {
		return nil, err
	}
	return result.Model, nil
}

func flattenHorizondbClusterPool(model *pools.Pool) ([]HorizondbClusterPool, error) {
	if model == nil {
		return []HorizondbClusterPool{}, nil
	}

	id, err := pools.ParsePoolIDInsensitively(pointer.From(model.Id))
	if err != nil {
		return []HorizondbClusterPool{}, err
	}

	props := model.Properties
	if props == nil {
		return []HorizondbClusterPool{}, nil
	}

	return []HorizondbClusterPool{{
		Name:                 id.PoolName,
		Id:                   id.ID(),
		NodeCount:            props.NodeCount,
		HighAvailabilityMode: string(props.HighAvailability.Mode),
		Compute: []HorizondbClusterPoolCompute{{
			Processor:  string(props.ComputeModel.ComputeProcessor),
			Generation: props.ComputeModel.ComputeGeneration,
			VCores:     props.ComputeModel.VCores,
		}},
	}}, nil
}

func createHorizondbClusterPasswordAuthentication(ctx context.Context, metadata sdk.ResourceMetaData, id clusters.ClusterId, state HorizondbClusterResourceModel) error {
	woPassword, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "administrator_password_wo", cty.String)
	if err != nil {
		return err
	}

	properties := passwordauthentications.PasswordAuthenticationProperties{
		State:                      passwordauthentications.AuthenticationStateEnabled,
		AdministratorLogin:         pointer.To(state.AdministratorLogin),
		AdministratorLoginPassword: pointer.To(state.AdministratorPassword),
	}

	if !woPassword.IsNull() {
		properties.AdministratorLoginPassword = pointer.To(woPassword.AsString())
	}

	authenticationId := passwordauthentications.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)

	payload := passwordauthentications.PasswordAuthentication{Properties: &properties}

	if err := metadata.Client.HorizonDB.PasswordAuthenticationsClient.CreateOrReplaceThenPoll(ctx, authenticationId, payload); err != nil {
		return fmt.Errorf("configuring password authentication for %s: %+v", id, err)
	}
	return nil
}

func readHorizondbClusterAdministratorLogin(ctx context.Context, client *passwordauthentications.PasswordAuthenticationsClient, id clusters.ClusterId) (string, error) {
	authenticationId := passwordauthentications.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)

	resp, err := client.Get(ctx, authenticationId)
	if err != nil {
		return "", fmt.Errorf("retrieving password authentication for %s: %+v", id, err)
	}

	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retrieving password authentication for %s: properties are missing", id)
	}
	return pointer.From(resp.Model.Properties.AdministratorLogin), nil
}
