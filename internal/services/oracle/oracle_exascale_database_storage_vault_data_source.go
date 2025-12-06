// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ExascaleDatabaseStorageVaultDataSource struct{}

type ExascaleDatabaseStorageVaultDataModel struct {
	Location          string       `tfschema:"location"`
	Name              string       `tfschema:"name"`
	ResourceGroupName string       `tfschema:"resource_group_name"`
	Zones             zones.Schema `tfschema:"zones"`

	AdditionalFlashCachePercentage int64                                 `tfschema:"additional_flash_cache_percentage"`
	Description                    string                                `tfschema:"description"`
	DisplayName                    string                                `tfschema:"display_name"`
	HighCapacityDatabaseStorage    []ExascaleDatabaseStorageDetailsModel `tfschema:"high_capacity_database_storage"`
	LifecycleDetails               string                                `tfschema:"lifecycle_details"`
	LifecycleState                 string                                `tfschema:"lifecycle_state"`
	Ocid                           string                                `tfschema:"ocid"`
	OciUrl                         string                                `tfschema:"oci_url"`
	TimeZone                       string                                `tfschema:"time_zone"`
	VirtualMachineClusterCount     int64                                 `tfschema:"virtual_machine_cluster_count"`
}

func (d ExascaleDatabaseStorageVaultDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseResourceName,
			),
		},
	}
}

func (d ExascaleDatabaseStorageVaultDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"additional_flash_cache_percentage": {
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

func (d ExascaleDatabaseStorageVaultDataSource) ModelObject() interface{} {
	return &ExascaleDatabaseStorageVaultDataModel{}
}

func (d ExascaleDatabaseStorageVaultDataSource) ResourceType() string {
	return "azurerm_oracle_exascale_database_storage_vault"
}

func (d ExascaleDatabaseStorageVaultDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exascaledbstoragevaults.ValidateExascaleDbStorageVaultID
}

func (d ExascaleDatabaseStorageVaultDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ExascaleDatabaseStorageVaultDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := exascaledbstoragevaults.NewExascaleDbStorageVaultID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return err
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.AdditionalFlashCachePercentage = pointer.From(props.AdditionalFlashCacheInPercent)
					state.DisplayName = props.DisplayName
					state.Description = pointer.From(props.Description)
					state.HighCapacityDatabaseStorage = FlattenHighCapacityDatabaseStorage(props.HighCapacityDatabaseStorage)
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
