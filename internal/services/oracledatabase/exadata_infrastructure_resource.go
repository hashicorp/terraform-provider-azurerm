// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracledatabase/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"compute_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.ComputeCount,
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
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.StorageCount,
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
					"days_of_week": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.DaysOfWeek,
						},
					},

					"hours_of_day": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validate.HoursOfDay,
						},
					},

					"lead_time_in_weeks": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validate.LeadTimeInWeeks,
					},

					"months": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.Month,
						},
					},

					"patching_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validate.PatchingMode,
					},

					"preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validate.Preference,
					},

					"weeks_of_month": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validate.WeeksOfMonth,
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequired(),
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
		Timeout: 120 * time.Minute,
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
					CustomerContacts: pointer.To(ExpandCustomerContacts(model.CustomerContacts)),
				},
			}

			if len(model.MaintenanceWindow) > 0 {
				param.Properties.MaintenanceWindow = &cloudexadatainfrastructures.MaintenanceWindow{
					DaysOfWeek:      pointer.To(ExpandDayOfWeekTo(model.MaintenanceWindow[0].DaysOfWeek)),
					HoursOfDay:      pointer.To(model.MaintenanceWindow[0].HoursOfDay),
					LeadTimeInWeeks: pointer.To(model.MaintenanceWindow[0].LeadTimeInWeeks),
					Months:          pointer.To(ExpandMonths(model.MaintenanceWindow[0].Months)),
					PatchingMode:    pointer.To(cloudexadatainfrastructures.PatchingMode(model.MaintenanceWindow[0].PatchingMode)),
					Preference:      pointer.To(cloudexadatainfrastructures.Preference(model.MaintenanceWindow[0].Preference)),
					WeeksOfMonth:    pointer.To(model.MaintenanceWindow[0].WeeksOfMonth),
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

			prop := result.Model.Properties
			output := ExadataInfraResourceModel{
				CustomerContacts:  FlattenCustomerContacts(result.Model.Properties.CustomerContacts),
				Name:              pointer.ToString(result.Model.Name),
				Location:          result.Model.Location,
				Zones:             result.Model.Zones,
				ResourceGroupName: id.ResourceGroupName,
				Tags:              utils.FlattenPtrMapStringString(result.Model.Tags),
				ComputeCount:      pointer.From(prop.ComputeCount),
				DisplayName:       prop.DisplayName,
				StorageCount:      pointer.From(prop.StorageCount),
				Shape:             prop.Shape,
				MaintenanceWindow: FlattenMaintenanceWindow(prop.MaintenanceWindow),
			}

			return metadata.Encode(&output)
		},
	}
}

func (ExadataInfraResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
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

func ExpandCustomerContacts(customerContactsList []string) []cloudexadatainfrastructures.CustomerContact {
	var customerContacts []cloudexadatainfrastructures.CustomerContact
	for _, customerContact := range customerContactsList {
		customerContacts = append(customerContacts, cloudexadatainfrastructures.CustomerContact{
			Email: customerContact,
		})
	}
	return customerContacts
}

func ExpandDayOfWeekTo(daysOfWeek []string) []cloudexadatainfrastructures.DayOfWeek {
	var daysOfWeekConverted []cloudexadatainfrastructures.DayOfWeek
	for _, day := range daysOfWeek {
		daysOfWeekConverted = append(daysOfWeekConverted, cloudexadatainfrastructures.DayOfWeek{
			Name: cloudexadatainfrastructures.DayOfWeekName(day),
		})
	}
	return daysOfWeekConverted
}

func ExpandMonths(months []string) []cloudexadatainfrastructures.Month {
	var monthsConverted []cloudexadatainfrastructures.Month
	for _, month := range months {
		monthsConverted = append(monthsConverted, cloudexadatainfrastructures.Month{
			Name: cloudexadatainfrastructures.MonthName(month),
		})
	}
	return monthsConverted
}
