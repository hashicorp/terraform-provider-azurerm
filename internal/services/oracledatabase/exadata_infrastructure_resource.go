// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

var _ sdk.Resource = ExadataInfraResource{}

type ExadataInfraResource struct{}

type ExadataInfraResourceModel struct {
	// Azure
	Location          string                 `tfschema:"location"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
	Zones             zones.Schema           `tfschema:"zones"`

	// Required
	ComputeCount int64  `tfschema:"compute_count"`
	DisplayName  string `tfschema:"display_name"`
	Shape        string `tfschema:"shape"`
	StorageCount int64  `tfschema:"storage_count"`

	// Optional
	CustomerContacts  []string                 `tfschema:"customer_contacts"`
	MaintenanceWindow []MaintenanceWindowModel `tfschema:"maintenance_window"`
}

func (ExadataInfraResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Azure
		"location": commonschema.Location(),
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"tags":  commonschema.Tags(),
		"zones": commonschema.ZonesMultipleRequired(),

		// Required
		"compute_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"shape": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"storage_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		// Optional
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"custom_action_timeout_in_mins": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},
					"days_of_week": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"hours_of_day": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},
					"is_custom_action_timeout_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},
					"is_monthly_patching_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},
					"lead_time_in_weeks": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},
					"months": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"patching_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"preference": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"weeks_of_month": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},
				},
			},
		},
	}
}

func (ExadataInfraResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ExadataInfraResource) ModelObject() interface{} {
	return &ExadataInfraResource{}
}

func (ExadataInfraResource) ResourceType() string {
	return "azurerm_oracledatabase_exadata_infrastructure"
}

func (r ExadataInfraResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudExadataInfrastructures
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ExadataInfraResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := cloudexadatainfrastructures.CloudExadataInfrastructure{
				Name:     pointer.To(model.Name),
				Location: model.Location,
				Tags:     tags.Expand(model.Tags),
				Zones:    model.Zones,
				Properties: &cloudexadatainfrastructures.CloudExadataInfrastructureProperties{
					ComputeCount:     pointer.To(model.ComputeCount),
					DisplayName:      model.DisplayName,
					StorageCount:     pointer.To(model.StorageCount),
					Shape:            model.Shape,
					CustomerContacts: pointer.To(ConvertCustomerContactsToSDK(model.CustomerContacts)),
				},
			}

			if model.MaintenanceWindow != nil && len(model.MaintenanceWindow) > 0 {
				param.Properties.MaintenanceWindow = &cloudexadatainfrastructures.MaintenanceWindow{
					CustomActionTimeoutInMins:    pointer.To(model.MaintenanceWindow[0].CustomActionTimeoutInMins),
					DaysOfWeek:                   pointer.To(ConvertDayOfWeekToSDK(model.MaintenanceWindow[0].DaysOfWeek)),
					HoursOfDay:                   pointer.To(model.MaintenanceWindow[0].HoursOfDay),
					IsCustomActionTimeoutEnabled: pointer.To(model.MaintenanceWindow[0].IsCustomActionTimeoutEnabled),
					IsMonthlyPatchingEnabled:     pointer.To(model.MaintenanceWindow[0].IsMonthlyPatchingEnabled),
					LeadTimeInWeeks:              pointer.To(model.MaintenanceWindow[0].LeadTimeInWeeks),
					Months:                       pointer.To(ConvertMonthsToSDK(model.MaintenanceWindow[0].Months)),
					PatchingMode:                 pointer.To(cloudexadatainfrastructures.PatchingMode(model.MaintenanceWindow[0].PatchingMode)),
					Preference:                   pointer.To(cloudexadatainfrastructures.Preference(model.MaintenanceWindow[0].Preference)),
					WeeksOfMonth:                 pointer.To(model.MaintenanceWindow[0].WeeksOfMonth),
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ExadataInfraResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudExadataInfrastructures
			id, err := cloudexadatainfrastructures.ParseCloudExadataInfrastructureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExadataInfraResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving exists when updating: +%v", *id)
			}
			if existing.Model == nil && existing.Model.Properties == nil {
				return fmt.Errorf("retrieving as nil when updating for %v", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				update := &cloudexadatainfrastructures.CloudExadataInfrastructureUpdate{
					Tags: tags.Expand(model.Tags),
				}
				err = client.UpdateThenPoll(ctx, *id, *update)
				if err != nil {
					return fmt.Errorf("updating %s: %v", id, err)
				}
			} else if metadata.ResourceData.HasChangesExcept("tags") {
				return fmt.Errorf("only `tags` currently support updates")
			}
			return nil
		},
	}
}

func (ExadataInfraResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := cloudexadatainfrastructures.ParseCloudExadataInfrastructureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudExadataInfrastructures
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}
			var output ExadataInfraResourceModel

			output.CustomerContacts = ConvertCustomerContactsToInternalModel(result.Model.Properties.CustomerContacts)
			output.Name = pointer.ToString(result.Model.Name)
			output.Location = result.Model.Location
			output.Zones = result.Model.Zones
			output.ResourceGroupName = id.ResourceGroupName
			output.Tags = utils.FlattenPtrMapStringString(result.Model.Tags)
			prop := result.Model.Properties
			output.ComputeCount = pointer.From(prop.ComputeCount)
			output.DisplayName = prop.DisplayName
			output.StorageCount = pointer.From(prop.StorageCount)
			output.Shape = prop.Shape
			output.MaintenanceWindow = ConvertMaintenanceWindowToInternalModel(prop.MaintenanceWindow)

			return metadata.Encode(&output)
		},
	}
}

func (ExadataInfraResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudExadataInfrastructures

			id, err := cloudexadatainfrastructures.ParseCloudExadataInfrastructureID(metadata.ResourceData.Id())
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

func (ExadataInfraResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID
}
