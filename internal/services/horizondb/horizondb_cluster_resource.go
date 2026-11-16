// Copyright IBM Corp. 2014, 2025
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

//go:generate go run ../../tools/generator-tests resourceidentity
type HorizondbClusterResource struct{}

type HorizondbClusterResourceModel struct {
	Name                           string                       `tfschema:"name"`
	ResourceGroupName              string                       `tfschema:"resource_group_name"`
	Location                       string                       `tfschema:"location"`
	PostgresqlVersion              string                       `tfschema:"postgresql_version"`
	ReadWritePool                  []HorizondbClusterPool       `tfschema:"read_write_pool"`
	AdministratorLogin             string                       `tfschema:"administrator_login"`
	AdministratorPassword          string                       `tfschema:"administrator_password"`
	AdministratorPasswordWOVersion int64                        `tfschema:"administrator_password_wo_version"`
	BackupRetentionDays            int64                        `tfschema:"backup_retention_days"`
	CreateMode                     string                       `tfschema:"create_mode"`
	GeoRedundantBackupEnabled      bool                         `tfschema:"geo_redundant_backup_enabled"`
	Identity                       []identity.ModelUserAssigned `tfschema:"identity"`
	PublicNetworkAccessEnabled     bool                         `tfschema:"public_network_access_enabled"`
	Restore                        []HorizondbClusterRestore    `tfschema:"restore"`
	Tags                           map[string]string            `tfschema:"tags"`
}

type HorizondbClusterRestore struct {
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
	Type       string `tfschema:"type"`
	Processor  string `tfschema:"processor"`
	Generation int64  `tfschema:"generation"`
	VCores     int64  `tfschema:"vcores"`
}

var (
	_ sdk.ResourceWithUpdate        = HorizondbClusterResource{}
	_ sdk.ResourceWithIdentity      = HorizondbClusterResource{}
	_ sdk.ResourceWithCustomizeDiff = HorizondbClusterResource{}
)

func (HorizondbClusterResource) ModelObject() interface{} {
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
		"location":            commonschema.Location(),
		"administrator_login": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 63),
		},
		"read_write_pool": horizondbClusterPoolSchema(),
		"administrator_password": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"administrator_password_wo"},
			ExactlyOneOf:  []string{"administrator_password", "administrator_password_wo"},
			ValidateFunc:  validation.StringLenBetween(8, 128),
		},
		"administrator_password_wo": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			ConflictsWith: []string{"administrator_password"},
			ExactlyOneOf:  []string{"administrator_password", "administrator_password_wo"},
			RequiredWith:  []string{"administrator_password_wo_version"},
			ValidateFunc:  validation.StringLenBetween(8, 128),
		},
		"administrator_password_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_password_wo"},
			ValidateFunc: validation.IntAtLeast(1),
		},
		"backup_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 35),
		},
		"create_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(clusters.ClusterCreateModeCreate),
			ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForClusterCreateMode(), false),
		},
		"geo_redundant_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
		"identity": commonschema.UserAssignedIdentityOptional(),
		"postgresql_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C - Restored clusters inherit their PostgreSQL version from the source backup.
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForPostgreSqlVersion(), false),
		},
		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"restore": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"point_in_time_utc": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
					"source_horizondb_cluster_id": commonschema.ResourceIDReferenceRequired(&clusters.ClusterId{}),
					"source_location": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						StateFunc:    location.StateFunc,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
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
			if err := validateHorizondbClusterConfig(diff.GetRawConfig()); err != nil {
				return err
			}
			if err := validateHorizondbClusterPoolConfig(diff.GetRawConfig()); err != nil {
				return err
			}
			if diff.Id() != "" && !diff.HasChanges("name", "resource_group_name", "location", "administrator_login", "create_mode", "postgresql_version", "restore", "geo_redundant_backup_enabled") && diff.NewValueKnown("read_write_pool.0.name") {
				old, new := diff.GetChange("read_write_pool.0.name")
				if err := validateHorizondbClusterPoolRename(old.(string), new.(string)); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func validateHorizondbClusterConfig(config cty.Value) error {
	if !config.IsKnown() || config.IsNull() {
		return nil
	}
	restore := config.GetAttr("restore")
	version := config.GetAttr("postgresql_version")
	createMode := config.GetAttr("create_mode")
	if createMode.IsKnown() {
		mode := string(clusters.ClusterCreateModeCreate)
		if !createMode.IsNull() {
			mode = createMode.AsString()
		}
		hasRestore := restore.IsKnown() && !restore.IsNull() && restore.LengthInt() > 0
		switch mode {
		case string(clusters.ClusterCreateModeCreate):
			if restore.IsKnown() && hasRestore {
				return errors.New("`restore` can only be specified when `create_mode` is `PointInTimeRestore`")
			}
			if version.IsNull() {
				return errors.New("`postgresql_version` is required when `create_mode` is `Create`")
			}
		case string(clusters.ClusterCreateModePointInTimeRestore):
			if restore.IsKnown() && !hasRestore {
				return errors.New("`restore` is required when `create_mode` is `PointInTimeRestore`")
			}
			if !version.IsNull() {
				return errors.New("`postgresql_version` must be omitted when `create_mode` is `PointInTimeRestore`")
			}
		}
	}
	return nil
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
			payload, err := expandHorizondbCluster(state)
			if err != nil {
				return err
			}
			if err := client.CreateOrReplaceCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			if err := createOrReplaceHorizondbClusterPool(ctx, metadata, id, state.ReadWritePool[0]); err != nil {
				return err
			}
			return reconcileHorizondbClusterPasswordAuthentication(ctx, metadata, id, state)
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
			var state HorizondbClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			oldName, newName := metadata.ResourceData.GetChange("read_write_pool.0.name")
			if err := validateHorizondbClusterPoolRename(oldName.(string), newName.(string)); err != nil {
				return err
			}
			if metadata.ResourceData.HasChanges("identity", "tags", "backup_retention_days", "public_network_access_enabled") {
				payload, err := expandHorizondbClusterUpdate(state)
				if err != nil {
					return err
				}
				if err := metadata.Client.HorizonDB.ClustersClient.CreateOrReplaceThenPoll(ctx, *id, payload); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}
			if metadata.ResourceData.HasChange("read_write_pool") {
				if err := createOrReplaceHorizondbClusterPool(ctx, metadata, *id, state.ReadWritePool[0]); err != nil {
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
			if err := metadata.Client.HorizonDB.ClustersClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func reconcileHorizondbClusterPasswordAuthentication(ctx context.Context, metadata sdk.ResourceMetaData, id clusters.ClusterId, state HorizondbClusterResourceModel) error {
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
	if resp.Model == nil || resp.Model.Properties == nil || pointer.From(resp.Model.Properties.AdministratorLogin) == "" {
		return "", fmt.Errorf("retrieving password authentication for %s: administrator login is missing", id)
	}
	return *resp.Model.Properties.AdministratorLogin, nil
}

func expandHorizondbCluster(state HorizondbClusterResourceModel) (clusters.Cluster, error) {
	geoRedundancy := clusters.GeoRedundancyStateDisabled
	if state.GeoRedundantBackupEnabled {
		geoRedundancy = clusters.GeoRedundancyStateEnabled
	}
	payload := clusters.Cluster{
		Location: location.Normalize(state.Location),
		Tags:     pointer.To(state.Tags),
		Properties: &clusters.ClusterProperties{
			CreationMode: pointer.ToEnum[clusters.ClusterCreateMode](state.CreateMode),
			Network:      &clusters.Network{PublicAccess: pointer.To(horizondbPublicNetworkAccess(state.PublicNetworkAccessEnabled))},
			Backup: &clusters.BackupProperties{
				GeoRedundancy: geoRedundancy,
				RetentionDays: pointer.To(state.BackupRetentionDays),
			},
		},
	}
	if len(state.Identity) > 0 {
		expanded, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
		if err != nil {
			return payload, fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = expanded
	}
	switch clusters.ClusterCreateMode(state.CreateMode) {
	case clusters.ClusterCreateModeCreate:
		if len(state.Restore) > 0 {
			return payload, errors.New("`restore` can only be specified when `create_mode` is `PointInTimeRestore`")
		}
		if state.PostgresqlVersion == "" {
			return payload, errors.New("`postgresql_version` is required when `create_mode` is `Create`")
		}
		payload.Properties.PostgreSqlVersion = pointer.ToEnum[clusters.PostgreSqlVersion](state.PostgresqlVersion)
	case clusters.ClusterCreateModePointInTimeRestore:
		if len(state.Restore) == 0 {
			return payload, errors.New("`restore` is required when `create_mode` is `PointInTimeRestore`")
		}
		restore := state.Restore[0]
		payload.Properties.CreationSourceClusterResourceId = pointer.To(restore.SourceHorizondbClusterId)
		payload.Properties.PointInTimeUtc = pointer.To(restore.PointInTimeUtc)
		if restore.SourceLocation != "" {
			payload.Properties.CreationSourceClusterLocation = pointer.To(location.Normalize(restore.SourceLocation))
		}
	default:
		return payload, fmt.Errorf("unsupported `create_mode` %q", state.CreateMode)
	}
	return payload, nil
}

func expandHorizondbClusterUpdate(state HorizondbClusterResourceModel) (clusters.Cluster, error) {
	expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
	if err != nil {
		return clusters.Cluster{}, fmt.Errorf("expanding `identity`: %+v", err)
	}
	geoRedundancy := clusters.GeoRedundancyStateDisabled
	if state.GeoRedundantBackupEnabled {
		geoRedundancy = clusters.GeoRedundancyStateEnabled
	}
	return clusters.Cluster{
		Location: location.Normalize(state.Location),
		Identity: expandedIdentity,
		Tags:     pointer.To(state.Tags),
		Properties: &clusters.ClusterProperties{
			PostgreSqlVersion: pointer.ToEnum[clusters.PostgreSqlVersion](state.PostgresqlVersion),
			Network:           &clusters.Network{PublicAccess: pointer.To(horizondbPublicNetworkAccess(state.PublicNetworkAccessEnabled))},
			Backup: &clusters.BackupProperties{
				GeoRedundancy: geoRedundancy,
				RetentionDays: pointer.To(state.BackupRetentionDays),
			},
		},
	}, nil
}

func (HorizondbClusterResource) flatten(metadata sdk.ResourceMetaData, id *clusters.ClusterId, model *clusters.Cluster, pool *pools.Pool, administratorLogin string) error {
	if model == nil || model.Properties == nil {
		return fmt.Errorf("retrieving %s: cluster properties are missing", id)
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
		BackupRetentionDays:            7,
		CreateMode:                     pointer.FromEnum(props.CreationMode),
		Tags:                           pointer.From(model.Tags),
	}
	if state.CreateMode == "" {
		state.CreateMode = string(clusters.ClusterCreateModeCreate)
	}
	if props.Network != nil {
		state.PublicNetworkAccessEnabled = pointer.From(props.Network.PublicAccess) == clusters.PublicNetworkAccessStateEnabled
	}
	if props.Backup != nil {
		if props.Backup.RetentionDays != nil {
			state.BackupRetentionDays = *props.Backup.RetentionDays
		}
		state.GeoRedundantBackupEnabled = props.Backup.GeoRedundancy == clusters.GeoRedundancyStateEnabled
	}
	flattenedIdentity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	state.Identity = pointer.From(flattenedIdentity)
	if pointer.From(props.CreationMode) == clusters.ClusterCreateModePointInTimeRestore {
		restore := HorizondbClusterRestore{
			SourceHorizondbClusterId: pointer.From(props.CreationSourceClusterResourceId),
			SourceLocation:           location.Normalize(pointer.From(props.CreationSourceClusterLocation)),
			PointInTimeUtc:           pointer.From(props.PointInTimeUtc),
		}
		// Azure can omit creation-time restore metadata from GET responses.
		if len(previous.Restore) > 0 {
			if restore.SourceHorizondbClusterId == "" {
				restore.SourceHorizondbClusterId = previous.Restore[0].SourceHorizondbClusterId
			}
			if restore.SourceLocation == "" {
				restore.SourceLocation = location.Normalize(previous.Restore[0].SourceLocation)
			}
			if restore.PointInTimeUtc == "" {
				restore.PointInTimeUtc = previous.Restore[0].PointInTimeUtc
			}
		}
		if restore.SourceHorizondbClusterId != "" || len(previous.Restore) > 0 {
			sourceId, err := clusters.ParseClusterIDInsensitively(restore.SourceHorizondbClusterId)
			if err != nil {
				return err
			}
			restore.SourceHorizondbClusterId = sourceId.ID()
			state.Restore = []HorizondbClusterRestore{restore}
		}
	}
	state.ReadWritePool, err = flattenHorizondbClusterPool(*id, pool)
	if err != nil {
		return err
	}
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func horizondbPublicNetworkAccess(enabled bool) clusters.PublicNetworkAccessState {
	if enabled {
		return clusters.PublicNetworkAccessStateEnabled
	}
	return clusters.PublicNetworkAccessStateDisabled
}

func horizondbClusterPoolSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
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
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pools.PoolComputeModelTypeProvisioned),
						ValidateFunc: validation.StringInSlice(pools.PossibleValuesForPoolComputeModelType(), false),
					},
					"processor": {
						Type: pluginsdk.TypeString, Required: true,
						ValidateFunc: validation.StringInSlice(pools.PossibleValuesForComputeProcessorType(), false),
					},
					"generation": {
						Type: pluginsdk.TypeInt, Required: true,
						ValidateFunc: validation.IntBetween(1, 2147483647),
					},
					"vcores": {
						Type: pluginsdk.TypeInt, Required: true,
						ValidateFunc: validation.IntBetween(2, 192),
					},
				}},
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
		}},
	}
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
		if (nodes == 1) != (mode.AsString() == string(pools.PoolHighAvailabilityModeDisabled)) {
			return errors.New("`read_write_pool.high_availability_mode` must be `Disabled` for one node, and `AnyZone` or `ZoneRedundant` for two or more nodes")
		}
	}
	return nil
}

func validateHorizondbClusterPoolRename(old, new string) error {
	if old != "" && new != "" && old != new {
		return errors.New("changing `read_write_pool.name` in place is not supported; replace the cluster to use a different pool name")
	}
	return nil
}

func readHorizondbClusterPool(ctx context.Context, metadata sdk.ResourceMetaData, id clusters.ClusterId, cluster *clusters.Cluster) (*pools.Pool, error) {
	if cluster == nil || cluster.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: cluster properties are missing", id)
	}
	designatedId := pointer.From(cluster.Properties.ReadWritePool.DesignatedPoolResourceId)
	if designatedId == "" {
		return nil, fmt.Errorf("retrieving %s: `readWritePool.designatedPoolResourceId` is missing", id)
	}
	poolId, err := pools.ParsePoolIDInsensitively(designatedId)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(poolId.SubscriptionId, id.SubscriptionId) || !strings.EqualFold(poolId.ResourceGroupName, id.ResourceGroupName) || !strings.EqualFold(poolId.ClusterName, id.ClusterName) {
		return nil, fmt.Errorf("designated read/write pool %s does not belong to %s", poolId, id)
	}
	result, err := metadata.Client.HorizonDB.PoolsClient.Get(ctx, *poolId)
	if err != nil {
		return nil, fmt.Errorf("retrieving designated read/write pool %s: %+v", poolId, err)
	}
	if result.Model == nil {
		return nil, fmt.Errorf("retrieving %s: pool model is missing", poolId)
	}
	returnedId, err := pools.ParsePoolIDInsensitively(pointer.From(result.Model.Id))
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(returnedId.ID(), poolId.ID()) {
		return nil, fmt.Errorf("retrieving %s: response pool ID does not match the designated read/write pool", poolId)
	}
	if _, err := flattenHorizondbClusterPool(id, result.Model); err != nil {
		return nil, err
	}
	return result.Model, nil
}

func flattenHorizondbClusterPool(clusterId clusters.ClusterId, model *pools.Pool) ([]HorizondbClusterPool, error) {
	if model == nil {
		return []HorizondbClusterPool{}, nil
	}
	id, err := pools.ParsePoolIDInsensitively(pointer.From(model.Id))
	if err != nil {
		return []HorizondbClusterPool{}, err
	}
	if !strings.EqualFold(id.SubscriptionId, clusterId.SubscriptionId) || !strings.EqualFold(id.ResourceGroupName, clusterId.ResourceGroupName) || !strings.EqualFold(id.ClusterName, clusterId.ClusterName) {
		return []HorizondbClusterPool{}, fmt.Errorf("pool %s does not belong to %s", id, clusterId)
	}
	props := model.Properties
	if props == nil || props.ComputeModel.ComputeProcessor == "" || props.ComputeModel.ComputeGeneration < 1 || props.ComputeModel.VCores < 2 || props.NodeCount < 1 || props.HighAvailability.Mode == "" {
		return []HorizondbClusterPool{}, fmt.Errorf("retrieving %s: required pool configuration is missing", id)
	}
	if props.ComputeModel.Type != pools.PoolComputeModelTypeProvisioned {
		return []HorizondbClusterPool{}, fmt.Errorf("retrieving %s: unsupported compute model type %q", id, props.ComputeModel.Type)
	}
	if props.ScheduledMaintenance != nil {
		maintenance := props.ScheduledMaintenance.ScheduledMaintenanceProperties().MaintenanceManagementType
		if maintenance != pools.ScheduledMaintenanceManagementTypeSystem {
			return []HorizondbClusterPool{}, fmt.Errorf("retrieving %s: unsupported maintenance management type %q", id, maintenance)
		}
	}
	return []HorizondbClusterPool{{
		Name: id.PoolName, Id: id.ID(), NodeCount: props.NodeCount,
		HighAvailabilityMode: string(props.HighAvailability.Mode),
		Compute: []HorizondbClusterPoolCompute{{
			Type:      string(props.ComputeModel.Type),
			Processor: string(props.ComputeModel.ComputeProcessor), Generation: props.ComputeModel.ComputeGeneration,
			VCores: props.ComputeModel.VCores,
		}},
	}}, nil
}

func createOrReplaceHorizondbClusterPool(ctx context.Context, metadata sdk.ResourceMetaData, clusterId clusters.ClusterId, state HorizondbClusterPool) error {
	id := pools.NewPoolID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ClusterName, state.Name)
	lockId := strings.ToLower(id.ID())
	locks.ByID(lockId)
	defer locks.UnlockByID(lockId)

	result, err := metadata.Client.HorizonDB.ClustersClient.Get(ctx, clusterId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", clusterId, err)
	}
	if result.Model == nil || result.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: cluster properties are missing", clusterId)
	}
	if pointer.From(result.Model.Properties.ReadWritePool.DesignatedPoolResourceId) != "" || !metadata.ResourceData.IsNewResource() {
		existing, err := readHorizondbClusterPool(ctx, metadata, clusterId, result.Model)
		if err != nil {
			return err
		}
		existingId, err := pools.ParsePoolIDInsensitively(pointer.From(existing.Id))
		if err != nil {
			return err
		}
		if !strings.EqualFold(existingId.ID(), id.ID()) {
			return fmt.Errorf("creating or replacing %s: remote pool changed; refresh state before retrying", id)
		}
	}
	payload := pools.Pool{
		Properties: &pools.PoolProperties{
			ComputeModel: expandHorizondbClusterPoolCompute(state.Compute[0]),
			NodeCount:    state.NodeCount, HighAvailability: pools.PoolHighAvailability{Mode: pools.PoolHighAvailabilityMode(state.HighAvailabilityMode)},
			ScheduledMaintenance: pools.SystemScheduledMaintenanceProperties{},
		},
	}
	if err := metadata.Client.HorizonDB.PoolsClient.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating or replacing %s: %+v", id, err)
	}
	return nil
}

func expandHorizondbClusterPoolCompute(input HorizondbClusterPoolCompute) pools.PoolComputeModel {
	return pools.PoolComputeModel{
		Type:              pools.PoolComputeModelType(input.Type),
		ComputeProcessor:  pools.ComputeProcessorType(input.Processor),
		ComputeGeneration: input.Generation, VCores: input.VCores,
	}
}
