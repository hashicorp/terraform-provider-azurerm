// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExascaleDbStorageVaultDataSource struct{}

type ExascaleDbStorageVaultDataModel struct {
	Location          string       `tfschema:"location"`
	Name              string       `tfschema:"name"`
	ResourceGroupName string       `tfschema:"resource_group_name"`
	Zones             zones.Schema `tfschema:"zones"`

	// ExascaleDbStorageVaultProperties
	AdditionalFlashCacheInPercent int64                           `tfschema:"additional_flash_cache_in_percent"`
	Description                   string                          `tfschema:"description"`
	DisplayName                   string                          `tfschema:"display_name"`
	HighCapacityDatabaseStorage   []ExascaleDbStorageDetailsModel `tfschema:"high_capacity_database_storage"`
	LifecycleDetails              string                          `tfschema:"lifecycle_details"`
	LifecycleState                string                          `tfschema:"lifecycle_state"`
	Ocid                          string                          `tfschema:"ocid"`
	OciUrl                        string                          `tfschema:"oci_url"`
	TimeZone                      string                          `tfschema:"time_zone"`
	VmClusterCount                int64                           `tfschema:"vm_cluster_count"`
}

type ExascaleDbStorageDetailsModel struct {
	AvailableSizeInGbs int64 `tfschema:"available_size_in_gbs"`
	TotalSizeInGbs     int64 `tfschema:"total_size_in_gbs"`
}

func (d ExascaleDbStorageVaultDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
		"zones":               commonschema.ZonesMultipleOptional(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageVaultName,
		},
	}
}

func (d ExascaleDbStorageVaultDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// ExascaleDbStorageVaultProperties
		"additional_flash_cache_in_percent": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"high_capacity_database_storage": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"available_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"total_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_details": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vm_cluster_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oci_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d ExascaleDbStorageVaultDataSource) ModelObject() interface{} {
	return &ExascaleDbStorageVaultDataModel{}
}

func (d ExascaleDbStorageVaultDataSource) ResourceType() string {
	return "azurerm_oracle_exascale_db_storage_vault"
}

func (d ExascaleDbStorageVaultDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exascaledbstoragevaults.ValidateExascaleDbStorageVaultID
}

func (d ExascaleDbStorageVaultDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient25.ExascaleDbStorageVaults
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ExascaleDbStorageVaultDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := exascaledbstoragevaults.NewExascaleDbStorageVaultID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {

					state.AdditionalFlashCacheInPercent = pointer.From(props.AdditionalFlashCacheInPercent)
					state.Description = pointer.From(props.Description)
					state.DisplayName = props.DisplayName
					state.HighCapacityDatabaseStorage = ExpandHighCapacityDatabaseStorage(props.HighCapacityDatabaseStorage)
					state.TimeZone = pointer.From(props.TimeZone)
					state.LifecycleState = string(*props.LifecycleState)
					state.LifecycleDetails = pointer.From(props.LifecycleDetails)
					state.VmClusterCount = pointer.From(props.VMClusterCount)
					state.Ocid = pointer.From(props.Ocid)
					state.OciUrl = pointer.From(props.OciURL)

				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func ExpandHighCapacityDatabaseStorage(input *exascaledbstoragevaults.ExascaleDbStorageDetails) []ExascaleDbStorageDetailsModel {
	output := make([]ExascaleDbStorageDetailsModel, 0)
	if input != nil {
		return append(output, ExascaleDbStorageDetailsModel{
			AvailableSizeInGbs: pointer.From(input.AvailableSizeInGbs),
			TotalSizeInGbs:     pointer.From(input.TotalSizeInGbs),
		})
	}
	return output
}
