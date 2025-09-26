// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
	VirtualMachineClusterCount    int64                           `tfschema:"virtual_machine_cluster_count"`
}

type ExascaleDbStorageDetailsModel struct {
	AvailableSizeInGb int64 `tfschema:"available_size_in_gb"`
	TotalSizeInGb     int64 `tfschema:"total_size_in_gb"`
}

func (d ExascaleDbStorageVaultDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_]`), "Name must start with a letter or underscore (_)"),
				validation.StringDoesNotContainAny("--"),
			),
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
					"available_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"total_size_in_gb": {
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

		"lifecycle_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_details": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"virtual_machine_cluster_count": {
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

		"zones": commonschema.ZonesMultipleComputed(),
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
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
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
					state.HighCapacityDatabaseStorage = flattenHighCapacityDatabaseStorage(props.HighCapacityDatabaseStorage)
					state.TimeZone = pointer.From(props.TimeZone)
					state.LifecycleState = pointer.FromEnum(props.LifecycleState)
					state.LifecycleDetails = pointer.From(props.LifecycleDetails)
					state.VirtualMachineClusterCount = pointer.From(props.VMClusterCount)
					state.Ocid = pointer.From(props.Ocid)
					state.OciUrl = pointer.From(props.OciURL)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenHighCapacityDatabaseStorage(input *exascaledbstoragevaults.ExascaleDbStorageDetails) []ExascaleDbStorageDetailsModel {
	output := make([]ExascaleDbStorageDetailsModel, 0)
	if input != nil {
		return append(output, ExascaleDbStorageDetailsModel{
			AvailableSizeInGb: pointer.From(input.AvailableSizeInGbs),
			TotalSizeInGb:     pointer.From(input.TotalSizeInGbs),
		})
	}
	return output
}
