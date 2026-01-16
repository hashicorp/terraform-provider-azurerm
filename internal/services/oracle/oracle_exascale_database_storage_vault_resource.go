// Copyright IBM Corp. 2014, 2025
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

var _ sdk.Resource = ExascaleDatabaseStorageVaultResource{}

type ExascaleDatabaseStorageVaultResource struct{}

type ExascaleDatabaseStorageVaultResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	AdditionalFlashCachePercentage int64                                 `tfschema:"additional_flash_cache_percentage"`
	Description                    string                                `tfschema:"description"`
	DisplayName                    string                                `tfschema:"display_name"`
	HighCapacityDatabaseStorage    []ExascaleDatabaseStorageDetailsModel `tfschema:"high_capacity_database_storage"`

	TimeZone string `tfschema:"time_zone"`
}

func (ExascaleDatabaseStorageVaultResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseResourceName,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"additional_flash_cache_percentage": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"description": {
			Type: pluginsdk.TypeString,
			// Note: O+C API use display_name value if omitted
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseResourceName,
			),
		},

		"high_capacity_database_storage": {
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},
					"available_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "UTC",
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (ExascaleDatabaseStorageVaultResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ExascaleDatabaseStorageVaultResource) ModelObject() interface{} {
	return &ExascaleDatabaseStorageVaultResource{}
}

func (ExascaleDatabaseStorageVaultResource) ResourceType() string {
	return "azurerm_oracle_exascale_database_storage_vault"
}

func (r ExascaleDatabaseStorageVaultResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ExascaleDatabaseStorageVaultResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := exascaledbstoragevaults.NewExascaleDbStorageVaultID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := exascaledbstoragevaults.ExascaleDbStorageVault{
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Zones:    pointer.To(model.Zones),
				Tags:     pointer.To(model.Tags),
				Properties: &exascaledbstoragevaults.ExascaleDbStorageVaultProperties{
					DisplayName:                   model.DisplayName,
					AdditionalFlashCacheInPercent: pointer.To(model.AdditionalFlashCachePercentage),
					TimeZone:                      pointer.To(model.TimeZone),
					HighCapacityDatabaseStorageInput: exascaledbstoragevaults.ExascaleDbStorageInputDetails{
						TotalSizeInGbs: model.HighCapacityDatabaseStorage[0].TotalSizeInGb,
					},
				},
			}

			if model.Description != "" {
				param.Properties.Description = pointer.To(model.Description)
			}

			if err := client.CreateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ExascaleDatabaseStorageVaultResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExascaleDatabaseStorageVaultResourceModel
			if err = metadata.Decode(&model); err != nil {
				return err
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			update := &exascaledbstoragevaults.ExascaleDbStorageVaultTagsUpdate{
				Tags: pointer.To(model.Tags),
			}

			err = client.UpdateThenPoll(ctx, *id, *update)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (ExascaleDatabaseStorageVaultResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ExascaleDatabaseStorageVaultResourceModel{
				Name:              id.ExascaleDbStorageVaultName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.AdditionalFlashCachePercentage = pointer.From(props.AdditionalFlashCacheInPercent)
					state.Description = pointer.From(props.Description)
					state.HighCapacityDatabaseStorage = FlattenHighCapacityDatabaseStorage(props.HighCapacityDatabaseStorage)
					state.TimeZone = pointer.From(props.TimeZone)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (ExascaleDatabaseStorageVaultResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults

			id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (ExascaleDatabaseStorageVaultResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exascaledbstoragevaults.ValidateExascaleDbStorageVaultID
}
