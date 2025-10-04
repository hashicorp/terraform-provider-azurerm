// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbstoragevaults"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ExascaleDbStorageVaultResource{}

type ExascaleDbStorageVaultResource struct{}

type ExascaleDbStorageVaultResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	// Required
	AdditionalFlashCacheInPercent int64                                `tfschema:"additional_flash_cache_in_percent"`
	Description                   string                               `tfschema:"description"`
	DisplayName                   string                               `tfschema:"display_name"`
	HighCapacityDatabaseStorage   []ExascaleDbStorageInputDetailsModel `tfschema:"high_capacity_database_storage"`

	// Optional
	TimeZone string `tfschema:"time_zone"`
}

type ExascaleDbStorageInputDetailsModel struct {
	TotalSizeInGb int64 `tfschema:"total_size_in_gb"`
}

func (ExascaleDbStorageVaultResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_]`), "Name must start with a letter or underscore (_)"),
				validation.StringDoesNotContainAny("--"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"additional_flash_cache_in_percent": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_]`), "Display Name must start with a letter or underscore (_)"),
				validation.StringDoesNotContainAny("--"),
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

func (ExascaleDbStorageVaultResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ExascaleDbStorageVaultResource) ModelObject() interface{} {
	return &ExascaleDbStorageVaultResource{}
}

func (ExascaleDbStorageVaultResource) ResourceType() string {
	return "azurerm_oracle_exascale_db_storage_vault"
}

func (r ExascaleDbStorageVaultResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ExascaleDbStorageVaultResourceModel
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
					AdditionalFlashCacheInPercent: pointer.To(model.AdditionalFlashCacheInPercent),
					Description:                   pointer.To(model.Description),
					TimeZone:                      pointer.To(model.TimeZone),
					HighCapacityDatabaseStorageInput: exascaledbstoragevaults.ExascaleDbStorageInputDetails{
						TotalSizeInGbs: model.HighCapacityDatabaseStorage[0].TotalSizeInGb,
					},
				},
			}

			if err := client.CreateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ExascaleDbStorageVaultResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExascaleDbStorageVaultResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			update := &exascaledbstoragevaults.ExascaleDbStorageVaultTagsUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(model.Tags)
			}

			err = client.UpdateThenPoll(ctx, *id, *update)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (ExascaleDbStorageVaultResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			client := metadata.Client.Oracle.OracleClient.ExascaleDbStorageVaults
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ExascaleDbStorageVaultResourceModel{
				Name:              id.ExascaleDbStorageVaultName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.Location = model.Location
					state.Tags = pointer.From(model.Tags)
					state.Zones = pointer.From(model.Zones)
					state.DisplayName = props.DisplayName
					state.AdditionalFlashCacheInPercent = pointer.From(props.AdditionalFlashCacheInPercent)
					state.Description = pointer.From(props.Description)
					state.HighCapacityDatabaseStorage = flattenHighCapacityDatabaseStorageInputInterface(props.HighCapacityDatabaseStorageInput)
					state.TimeZone = pointer.From(props.TimeZone)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (ExascaleDbStorageVaultResource) Delete() sdk.ResourceFunc {
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

func (ExascaleDbStorageVaultResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exascaledbstoragevaults.ValidateExascaleDbStorageVaultID
}

func flattenHighCapacityDatabaseStorageInputInterface(input exascaledbstoragevaults.ExascaleDbStorageInputDetails) []ExascaleDbStorageInputDetailsModel {
	output := make([]ExascaleDbStorageInputDetailsModel, 0)
	storageInput := ExascaleDbStorageInputDetailsModel{
		TotalSizeInGb: input.TotalSizeInGbs,
	}
	output = append(output, storageInput)
	return output
}
