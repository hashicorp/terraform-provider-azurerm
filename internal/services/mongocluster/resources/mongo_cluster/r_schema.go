// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
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
			Default:  string(mongoclusters.CreateModeDefault),
			// Confirmed with service team the 'Default', `GeoReplica` and `PointInTimeRestore` are the only accepted value currently, other values will be supported later.
			ValidateFunc: validation.StringInSlice([]string{
				string(mongoclusters.CreateModeDefault),
				string(mongoclusters.CreateModeGeoReplica),
				string(mongoclusters.CreateModePointInTimeRestore),
			}, false),
		},

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeKey),
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"data_api_mode_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// Although Swagger defines four identity types, the service API currently only supports `None` and `UserAssigned`. Service team has confirmed that they will support the other types in the future.
		"identity": commonschema.UserAssignedIdentityOptional(),

		"preview_features": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
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
						ForceNew:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},

					"source_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: mongoclusters.ValidateMongoClusterID,
					},
				},
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

		"authentication_methods": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			// NOTE: O+C `AuthConfig` is an object and its sub property `AllowedModes` has default value `NativeAuth` when `AuthConfig` isn't set in the tfconfig. So, `O+C` is required otherwise it will incur difference.
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(mongoclusters.PossibleValuesForAuthenticationMode(), false),
			},
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

		"cosmos_db_network_bypass_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
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
			ValidateFunc: validation.IntBetween(32, 32768),
		},

		"storage_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(mongoclusters.StorageTypePremiumSSD),
			ValidateFunc: validation.StringInSlice(mongoclusters.PossibleValuesForStorageType(), false),
			RequiredWith: []string{"storage_size_in_gb"},
		},

		"tags": commonschema.Tags(),

		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"5.0",
				"6.0",
				"7.0",
				"8.0",
			}, false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
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

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state MongoClusterResourceModel
			if err := metadata.DecodeDiff(&state); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			switch state.CreateMode {
			case string(mongoclusters.CreateModeDefault):
				if isNativeAuthRequired(metadata) && state.AdministratorUserName == "" {
					return errors.New("`administrator_username` is required when `authentication_methods` contains `NativeAuth` or is not configured")
				}

				if state.ComputeTier == "" {
					return fmt.Errorf("`compute_tier` is required when `create_mode` is `%s`", string(mongoclusters.CreateModeDefault))
				}

				if state.StorageSizeInGb == 0 {
					return fmt.Errorf("`storage_size_in_gb` is required when `create_mode` is `%s`", string(mongoclusters.CreateModeDefault))
				}

				if state.HighAvailabilityMode == "" {
					return fmt.Errorf("`high_availability_mode` is required when `create_mode` is `%s`", string(mongoclusters.CreateModeDefault))
				}

				if state.ShardCount == 0 {
					return fmt.Errorf("`shard_count` is required when `create_mode` is `%s`", string(mongoclusters.CreateModeDefault))
				}

				if state.Version == "" {
					return fmt.Errorf("`version` is required when `create_mode` is `%s`", string(mongoclusters.CreateModeDefault))
				}
			case string(mongoclusters.CreateModeGeoReplica):
				if state.SourceLocation == "" {
					return fmt.Errorf("`source_location` is required when `create_mode` is `GeoReplica`")
				}
			case string(mongoclusters.CreateModePointInTimeRestore):
				if len(state.Restore) == 0 {
					return fmt.Errorf("`restore` is required when `create_mode` is `PointInTimeRestore`")
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

			// Since the API doesn't return the value of create_mode, when importing `azurerm_mongo_cluster`, it will cause replacement.
			if oldVal, newVal := metadata.ResourceDiff.GetChange("create_mode"); oldVal.(string) != "" && oldVal.(string) != newVal.(string) {
				if err := metadata.ResourceDiff.ForceNew("create_mode"); err != nil {
					return err
				}
			}

			if !metadata.ResourceDiff.GetRawConfig().AsValueMap()["data_api_mode_enabled"].IsNull() && state.CreateMode != string(mongoclusters.CreateModeDefault) {
				return fmt.Errorf("`data_api_mode_enabled` can only be set when `create_mode` is `Default`")
			}

			if err := validateMongoClusterCosmosDBNetworkBypassEnabled(metadata, state); err != nil {
				return err
			}

			// Service team confirmed that `data_api_mode_enabled` can only be updated to `Enabled` after the cluster has been created
			if oldVal, newVal := metadata.ResourceDiff.GetChange("data_api_mode_enabled"); oldVal.(bool) && !newVal.(bool) && state.CreateMode == string(mongoclusters.CreateModeDefault) {
				if err := metadata.ResourceDiff.ForceNew("data_api_mode_enabled"); err != nil {
					return err
				}
			}

			// When identity is added or removed from the configuration, it should trigger ForceNew.
			// 1. When identity type is changed from `UserAssigned` to `None`, go-azure-sdk always sets `identity.userAssignedIdentities` to `nil` in the request payload.
			// But since service API does not allow explicitly setting `identity.userAssignedIdentities` to `nil` when the identity type is `None`,
			// otherwise it will throw a schema validation error, Terraform can only dynamically treat this change as forceNew.
			// 2. Service API will fail when identity type is changed from `None` to `UserAssigned`.
			if oldVal, newVal := metadata.ResourceDiff.GetChange("identity"); (len(oldVal.([]any)) > 0 && len(newVal.([]any)) == 0) || (len(oldVal.([]any)) == 0 && len(newVal.([]any)) > 0) {
				if err := metadata.ResourceDiff.ForceNew("identity"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
