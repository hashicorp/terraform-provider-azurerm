// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2024-07-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MongoClusterResource struct{}

var _ sdk.ResourceWithUpdate = MongoClusterResource{}

var _ sdk.ResourceWithCustomizeDiff = MongoClusterResource{}

type MongoClusterResourceModel struct {
	Name                  string                         `tfschema:"name"`
	ResourceGroupName     string                         `tfschema:"resource_group_name"`
	Location              string                         `tfschema:"location"`
	AdministratorUserName string                         `tfschema:"administrator_username"`
	AdministratorPassword string                         `tfschema:"administrator_password"`
	CreateMode            string                         `tfschema:"create_mode"`
	ShardCount            int64                          `tfschema:"shard_count"`
	SourceLocation        string                         `tfschema:"source_location"`
	SourceServerId        string                         `tfschema:"source_server_id"`
	ComputeTier           string                         `tfschema:"compute_tier"`
	HighAvailabilityMode  string                         `tfschema:"high_availability_mode"`
	PublicNetworkAccess   string                         `tfschema:"public_network_access"`
	PreviewFeatures       []string                       `tfschema:"preview_features"`
	StorageSizeInGb       int64                          `tfschema:"storage_size_in_gb"`
	ConnectionStrings     []MongoClusterConnectionString `tfschema:"connection_strings"`
	Tags                  map[string]string              `tfschema:"tags"`
	Version               string                         `tfschema:"version"`
}

type MongoClusterConnectionString struct {
	Value       string `tfschema:"value"`
	Description string `tfschema:"description"`
	Name        string `tfschema:"name"`
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

func (r MongoClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z\d]([-a-z\d]{1,38}[a-z\d])$`),
				"`name` must be between 3 and 40 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_username": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"administrator_password"},
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
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
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"shard_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			ForceNew:     true,
		},

		"source_location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			ValidateFunc:     validation.StringIsNotEmpty,
			RequiredWith:     []string{"source_server_id"},
		},

		"source_server_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: mongoclusters.ValidateMongoClusterID,
		},

		"administrator_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"administrator_username"},
		},

		"compute_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Free",
				"M10",
				"M20",
				"M25",
				"M30",
				"M40",
				"M50",
				"M60",
				"M80",
				"M200",
			}, false),
		},

		"high_availability_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				// Confirmed with service team the `SameZone` is currently not supported.
				string(mongoclusters.HighAvailabilityModeDisabled),
				string(mongoclusters.HighAvailabilityModeZoneRedundantPreferred),
			}, false),
		},

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(mongoclusters.PublicNetworkAccessEnabled),
			ValidateFunc: validation.StringInSlice(mongoclusters.PossibleValuesForPublicNetworkAccess(), false),
		},

		"storage_size_in_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(32, 16384),
		},

		"tags": commonschema.Tags(),

		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"5.0",
				"6.0",
				"7.0",
			}, false),
		},
	}
}

func (r MongoClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"connection_strings": {
			Type:      pluginsdk.TypeList,
			Sensitive: true,
			Computed:  true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"value": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
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
				Location:   location.Normalize(state.Location),
				Properties: &mongoclusters.MongoClusterProperties{},
			}

			if state.AdministratorUserName != "" {
				parameter.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorUserName),
					Password: pointer.To(state.AdministratorPassword),
				}
			}

			if state.CreateMode != "" {
				parameter.Properties.CreateMode = pointer.To(mongoclusters.CreateMode(state.CreateMode))
			}

			parameter.Properties.PreviewFeatures = expandPreviewFeatures(state.PreviewFeatures)

			if state.ShardCount != 0 {
				parameter.Properties.Sharding = &mongoclusters.ShardingProperties{
					ShardCount: pointer.To(state.ShardCount),
				}
			}

			if state.CreateMode == string(mongoclusters.CreateModeGeoReplica) {
				if state.SourceServerId == "" {
					return fmt.Errorf("`source_server_id` is required when `create_mode` is `GeoReplica`")
				}

				parameter.Properties.ReplicaParameters = &mongoclusters.MongoClusterReplicaParameters{
					SourceLocation:   state.SourceLocation,
					SourceResourceId: state.SourceServerId,
				}
			}

			if state.ComputeTier != "" {
				parameter.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
			}

			if state.HighAvailabilityMode != "" {
				parameter.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.To(mongoclusters.HighAvailabilityMode(state.HighAvailabilityMode)),
				}
			}

			parameter.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccess(state.PublicNetworkAccess))

			if state.StorageSizeInGb != 0 {
				parameter.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
				}
			}

			if state.Tags != nil {
				parameter.Tags = pointer.To(state.Tags)
			}

			if state.Version != "" {
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
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}
			payload := existing.Model

			metadata.Logger.Infof("updating %s", *id)

			// Set SystemData to nil as the API returns `The property '#/systemData' of type null did not match the following type: object in schema 25debcc2-6915-5536-9566-a2ecd765b755"}}` error.
			// https://github.com/Azure/azure-rest-api-specs/issues/31377 has been filed to track it.
			payload.SystemData = nil

			// upgrades involving Free or M25(Burstable) compute tier require first upgrading the compute tier, after which other configurations can be updated.
			if metadata.ResourceData.HasChange("compute_tier") {
				payload.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
				oldComputeTier, newComputeTier := metadata.ResourceData.GetChange("compute_tier")
				if (oldComputeTier == "Free" || oldComputeTier == "M25") && newComputeTier != "Free" && newComputeTier != "M25" {
					metadata.Logger.Infof("updating compute tier for %s", *id)
					if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
						return fmt.Errorf("updating %s: %+v", *id, err)
					}
				}
			}

			metadata.Logger.Infof("updating other configurations for %s", *id)
			if metadata.ResourceData.HasChange("administrator_password") {
				payload.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorUserName),
					Password: pointer.To(state.AdministratorPassword),
				}
			}

			if metadata.ResourceData.HasChange("high_availability_mode") {
				payload.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.To(mongoclusters.HighAvailabilityMode(state.HighAvailabilityMode)),
				}
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.To(mongoclusters.PublicNetworkAccess(state.PublicNetworkAccess))
			}

			if metadata.ResourceData.HasChange("storage_size_in_gb") {
				payload.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("version") {
				payload.Properties.ServerVersion = pointer.To(state.Version)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
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
				state.Location = location.Normalize(model.Location)

				if props := model.Properties; props != nil {
					// API doesn't return the value of administrator_password
					state.AdministratorPassword = metadata.ResourceData.Get("administrator_password").(string)

					// API doesn't return the value of create_mode, https://github.com/Azure/azure-rest-api-specs/issues/31266 has been filed to track it.
					state.CreateMode = metadata.ResourceData.Get("create_mode").(string)

					if v := props.Administrator; v != nil {
						state.AdministratorUserName = pointer.From(v.UserName)
					}

					if v := props.Replica; v != nil {
						// API doesn't return the value of source_location, https://github.com/Azure/azure-rest-api-specs/issues/31266 has been filed to track it.
						state.SourceLocation = metadata.ResourceData.Get("source_location").(string)
						if v.SourceResourceId != nil {
							id, err := mongoclusters.ParseMongoClusterID(pointer.From(v.SourceResourceId))
							if err != nil {
								return err
							}
							state.SourceServerId = id.ID()
						}
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
					state.PublicNetworkAccess = string(pointer.From(props.PublicNetworkAccess))

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

			csResp, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing connection strings for %s: %+v", *id, err)
			}
			if model := csResp.Model; model != nil {
				state.ConnectionStrings = flattenMongoClusterConnectionStrings(model.ConnectionStrings, state.AdministratorUserName, state.AdministratorPassword)
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
				if state.AdministratorUserName == "" {
					return fmt.Errorf("`administrator_username` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if state.ComputeTier == "" {
					return fmt.Errorf("`compute_tier` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if state.StorageSizeInGb == 0 {
					return fmt.Errorf("`storage_size_in_gb` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if state.HighAvailabilityMode == "" {
					return fmt.Errorf("`high_availability_mode` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if state.ShardCount == 0 {
					return fmt.Errorf("`shard_count` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}

				if state.Version == "" {
					return fmt.Errorf("`version` is required when `create_mode` is %s", string(mongoclusters.CreateModeDefault))
				}
			case string(mongoclusters.CreateModeGeoReplica):
				if state.SourceLocation == "" {
					return fmt.Errorf("`source_location` is required when `create_mode` is `GeoReplica`")
				}
			}

			if state.ComputeTier == "Free" || state.ComputeTier == "M25" {
				if state.HighAvailabilityMode == string(mongoclusters.HighAvailabilityModeZoneRedundantPreferred) {
					return fmt.Errorf("high Availability is not available with the `Free` or `M25` Compute Tier")
				}

				if state.ShardCount > 1 {
					return fmt.Errorf("the value of `shard_count` cannot exceed 1 for the `Free` or `M25` Compute Tier")
				}
			}

			if len(state.PreviewFeatures) > 0 {
				existing := make(map[string]bool)
				for _, str := range state.PreviewFeatures {
					if existing[str] {
						return fmt.Errorf("`preview_features` contains the duplicate value %q", str)
					}
					existing[str] = true
				}
			}

			return nil
		},
	}
}

func expandPreviewFeatures(input []string) *[]mongoclusters.PreviewFeature {
	if len(input) == 0 {
		return nil
	}

	result := make([]mongoclusters.PreviewFeature, 0)

	for _, v := range input {
		if v != "" {
			result = append(result, mongoclusters.PreviewFeature(v))
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

func flattenMongoClusterConnectionStrings(input *[]mongoclusters.ConnectionString, userName, userPassword string) []MongoClusterConnectionString {
	results := make([]MongoClusterConnectionString, 0)
	if input == nil {
		return results
	}
	for _, cs := range *input {
		value := pointer.From(cs.ConnectionString)
		// Password can be empty if it isn't available in the state file (e.g. during import).
		// In this case, we simply leave the placeholder unchanged.
		if userPassword != "" {
			value = regexp.MustCompile(`<user>:<password>`).ReplaceAllString(value, url.UserPassword(userName, userPassword).String())
		}

		results = append(results, MongoClusterConnectionString{
			Name:        pointer.From(cs.Name),
			Description: pointer.From(cs.Description),
			Value:       value,
		})
	}

	return results
}
