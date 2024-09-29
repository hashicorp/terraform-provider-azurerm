// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2024-07-01/mongoclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MongoClusterResource struct{}

var _ sdk.ResourceWithUpdate = MongoClusterResource{}

var _ sdk.ResourceWithCustomizeDiff = MongoClusterResource{}

type MongoClusterResourceModel struct {
	Name                       string            `tfschema:"name"`
	ResourceGroupName          string            `tfschema:"resource_group_name"`
	Location                   string            `tfschema:"location"`
	AdministratorLogin         string            `tfschema:"administrator_login"`
	AdministratorLoginPassword string            `tfschema:"administrator_login_password"`
	CreateMode                 string            `tfschema:"create_mode"`
	ShardCount                 int64             `tfschema:"shard_count"`
	SourceLocation             string            `tfschema:"source_location"`
	SourceServerId             string            `tfschema:"source_server_id"`
	ComputeTier                string            `tfschema:"compute_tier"`
	HighAvailabilityMode       string            `tfschema:"high_availability_mode"`
	PublicNetworkAccessEnabled bool              `tfschema:"public_network_access_enabled"`
	PreviewFeatures            []string          `tfschema:"preview_features"`
	StorageSizeInGb            int64             `tfschema:"storage_size_in_gb"`
	Tags                       map[string]string `tfschema:"tags"`
	Version                    string            `tfschema:"version"`
}

func (r MongoClusterResource) ModelObject() interface{} {
	return &MongoClusterResourceModel{}
}

func (r MongoClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mongoclusters.ValidateMongoClusterID
}

func (r MongoClusterResource) ResourceType() string {
	return "azurerm_mongo_cluster"
}

func (r MongoClusterResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     schema.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 40),
				validation.StringMatch(
					regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`),
					"The name contains only lowercase letters, numbers and hyphens.",
				),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"create_mode": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(mongoclusters.CreateModeDefault),
			// Confirmed with service team the 'Default' and `GeoReplica` are the only accepted value currently, other values will be supported later.
			ValidateFunc: validation.StringInSlice([]string{
				string(mongoclusters.CreateModeDefault),
				string(mongoclusters.CreateModeGeoReplica),
			}, false),
		},

		"preview_features": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(mongoclusters.PossibleValuesForPreviewFeature(), false),
			},
		},

		"shard_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			ForceNew:     true,
		},

		"source_location": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			ValidateFunc:     validation.StringIsNotEmpty,
			RequiredWith:     []string{"source_server_id"},
		},

		"source_server_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"administrator_login_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"compute_tier": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Free",
				"M25",
				"M30",
				"M40",
				"M50",
				"M60",
				"M80",
			}, false),
		},

		"high_availability_mode": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				// Confirmed with service team the `SameZone` is currently not supported.
				string(mongoclusters.HighAvailabilityModeDisabled),
				string(mongoclusters.HighAvailabilityModeZoneRedundantPreferred),
			}, false),
		},

		"public_network_access_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},

		"storage_size_in_gb": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(32, 16384),
		},

		"tags": commonschema.Tags(),

		"version": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"5.0",
				"6.0",
				"7.0",
			}, false),
		},
	}
}

func (r MongoClusterResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r MongoClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state MongoClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := mongoclusters.NewMongoClusterID(subscriptionId, state.ResourceGroupName, state.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameter := mongoclusters.MongoCluster{
				Location:   azure.NormalizeLocation(state.Location),
				Properties: &mongoclusters.MongoClusterProperties{},
			}

			if _, ok := metadata.ResourceData.GetOk("administrator_login"); ok {
				parameter.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorLogin),
					Password: pointer.To(state.AdministratorLoginPassword),
				}
			}

			if _, ok := metadata.ResourceData.GetOk("create_mode"); ok {
				parameter.Properties.CreateMode = pointer.To(mongoclusters.CreateMode(state.CreateMode))
			}

			if v, ok := metadata.ResourceData.GetOk("preview_features"); ok {
				parameter.Properties.PreviewFeatures = expandPreviewFeatures(v.([]interface{}))
			}

			if _, ok := metadata.ResourceData.GetOk("shard_count"); ok {
				parameter.Properties.Sharding = &mongoclusters.ShardingProperties{
					ShardCount: pointer.To(state.ShardCount),
				}
			}

			if state.CreateMode == string(mongoclusters.CreateModeGeoReplica) {
				parameter.Properties.ReplicaParameters = &mongoclusters.MongoClusterReplicaParameters{
					SourceLocation:   state.SourceLocation,
					SourceResourceId: state.SourceServerId,
				}
			}

			if _, ok := metadata.ResourceData.GetOk("compute_tier"); ok {
				parameter.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
			}

			if _, ok := metadata.ResourceData.GetOk("high_availability_mode"); ok {
				parameter.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.To(mongoclusters.HighAvailabilityMode(state.HighAvailabilityMode)),
				}
			}

			if state.PublicNetworkAccessEnabled {
				parameter.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessEnabled)
			} else {
				parameter.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessDisabled)
			}

			if _, ok := metadata.ResourceData.GetOk("storage_size_in_gb"); ok {
				parameter.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
				}
			}

			if _, ok := metadata.ResourceData.GetOk("tags"); ok {
				parameter.Tags = pointer.To(state.Tags)
			}

			if _, ok := metadata.ResourceData.GetOk("version"); ok {
				parameter.Properties.ServerVersion = pointer.To(state.Version)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameter); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MongoClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state...")
			var state MongoClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s for update: model was nil", *id)
			}

			metadata.Logger.Infof("updating %s", *id)

			model.SystemData = nil

			// upgrades involving Free or M25(Burstable) cluster tier require first upgrading the cluster tier, after which other configurations can be updated.
			if metadata.ResourceData.HasChange("compute_tier") {
				model.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
				oldComputeTier, newComputeTier := metadata.ResourceData.GetChange("compute_tier")
				if (oldComputeTier == "Free" || oldComputeTier == "M25") && newComputeTier != "Free" && newComputeTier != "M25" {
					metadata.Logger.Infof("updating cluster tier for %s", *id)
					if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
						return fmt.Errorf("updating %s: %+v", *id, err)
					}
				}
			}

			metadata.Logger.Infof("updating other configurations for %s", *id)
			if metadata.ResourceData.HasChange("administrator_login_password") {
				model.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorLogin),
					Password: pointer.To(state.AdministratorLoginPassword),
				}
			}

			if metadata.ResourceData.HasChange("high_availability_mode") {
				model.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.To(mongoclusters.HighAvailabilityMode(state.HighAvailabilityMode)),
				}
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if state.PublicNetworkAccessEnabled {
					model.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessEnabled)
				} else {
					model.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccessDisabled)
				}
			}

			if metadata.ResourceData.HasChange("storage_size_in_gb") {
				model.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("version") {
				model.Properties.ServerVersion = pointer.To(state.Version)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MongoClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MongoClusterResourceModel{
				Name:              id.MongoClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(&model.Location)

				if props := model.Properties; props != nil {
					// API doesn't return the value of administrator_login_password
					state.AdministratorLoginPassword = metadata.ResourceData.Get("administrator_login_password").(string)

					// API doesn't return the value of create_mode
					state.CreateMode = metadata.ResourceData.Get("create_mode").(string)

					if v := props.Administrator; v != nil {
						state.AdministratorLogin = pointer.From(v.UserName)
					}

					if v := props.Replica; v != nil {
						// API doesn't return the value of source_location
						state.SourceLocation = metadata.ResourceData.Get("source_location").(string)
						state.SourceServerId = pointer.From(v.SourceResourceId)
					}

					if v := props.Sharding; v != nil {
						state.ShardCount = pointer.From(v.ShardCount)
					}
					if v := props.Compute; v != nil {
						state.ComputeTier = pointer.From(v.Tier)
					}

					if v := props.HighAvailability; v != nil {
						state.HighAvailabilityMode = string(pointer.From(v.TargetMode))
					}
					state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == mongoclusters.PublicNetworkAccessEnabled

					if v := props.Storage; v != nil {
						state.StorageSizeInGb = pointer.From(v.SizeGb)
					}
					if v := props.PreviewFeatures; v != nil {
						state.PreviewFeatures = flattenMongoClusterPreviewFeatures(v)
					}
					state.Version = pointer.From(props.ServerVersion)
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MongoClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MongoClusterResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state MongoClusterResourceModel
			if err := metadata.DecodeDiff(&state); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			switch state.CreateMode {
			case string(mongoclusters.CreateModeDefault):
				if _, ok := metadata.ResourceDiff.GetOk("administrator_login"); !ok {
					return fmt.Errorf("`administrator_login` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("administrator_login_password"); !ok {
					return fmt.Errorf("`administrator_login_password` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("compute_tier"); !ok {
					return fmt.Errorf("`compute_tier` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("storage_size_in_gb"); !ok {
					return fmt.Errorf("`storage_size_in_gb` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("high_availability_mode"); !ok {
					return fmt.Errorf("`high_availability_mode` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("shard_count"); !ok {
					return fmt.Errorf("`shard_count` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if _, ok := metadata.ResourceDiff.GetOk("version"); !ok {
					return fmt.Errorf("`version` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}
			case string(mongoclusters.CreateModeGeoReplica):
				if state.SourceServerId == "" {
					return fmt.Errorf("`source_server_id` is required when `create_mode` is `GeoReplica`")
				}
				if state.SourceLocation == "" {
					return fmt.Errorf("`source_location` is required when `create_mode` is `GeoReplica`")
				}
			}

			if state.ComputeTier == "Free" || state.ComputeTier == "M25" {
				if state.HighAvailabilityMode == string(mongoclusters.HighAvailabilityModeZoneRedundantPreferred) {
					return fmt.Errorf("high Availability is not available with the 'Free' or 'M25' Cluster Tier")
				}

				if state.ShardCount > 1 {
					return fmt.Errorf("the value of `shard_count` cannot exceed 1 for the 'Free' or 'M25' Cluster Tier")
				}
			}

			return nil
		},
	}
}

func expandPreviewFeatures(input []interface{}) *[]mongoclusters.PreviewFeature {
	result := make([]mongoclusters.PreviewFeature, 0)

	for _, v := range input {
		if v != "" {
			result = append(result, mongoclusters.PreviewFeature(v.(string)))
		}
	}

	return &result
}

func flattenMongoClusterPreviewFeatures(input *[]mongoclusters.PreviewFeature) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}
