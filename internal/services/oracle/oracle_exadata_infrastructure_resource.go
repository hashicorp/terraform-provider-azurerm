// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = ExadataInfraResource{}

type ExadataInfraResource struct{}

type ExadataInfraResourceModel struct {
	// Azure
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ExadataName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"compute_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.ComputeCount,
			ForceNew:     true,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.CloudVMClusterName,
			ForceNew:     true,
		},

		"shape": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"storage_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.StorageCount,
			ForceNew:     true,
		},

		// Optional
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"days_of_week": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.DaysOfWeek,
						},
					},

					"hours_of_day": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validate.HoursOfDay,
						},
					},

					"lead_time_in_weeks": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.LeadTimeInWeeks,
					},

					"months": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.Month,
						},
					},

					"patching_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.PatchingMode,
					},

					"preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.Preference,
					},

					"weeks_of_month": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validate.WeeksOfMonth,
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (ExadataInfraResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ExadataInfraResource) ModelObject() interface{} {
	return &ExadataInfraResource{}
}

func (ExadataInfraResource) ResourceType() string {
	return "azurerm_oracle_exadata_infrastructure"
}

func (r ExadataInfraResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudExadataInfrastructures
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
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
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
			client := metadata.Client.Oracle.OracleClient.CloudExadataInfrastructures
			id, err := cloudexadatainfrastructures.ParseCloudExadataInfrastructureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExadataInfraResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: ", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				update := &cloudexadatainfrastructures.CloudExadataInfrastructureUpdate{
					Tags: pointer.To(model.Tags),
				}
				err = client.UpdateThenPoll(ctx, *id, *update)
				if err != nil {
					return fmt.Errorf("updating %s: %v", id, err)
				}
			}
			return nil
		},
	}
}

func (ExadataInfraResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudExadataInfrastructures

			id, err := cloudexadatainfrastructures.ParseCloudExadataInfrastructureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ExadataInfraResourceModel{
				Name:              id.CloudExadataInfrastructureName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Zones = model.Zones
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.CustomerContacts = FlattenCustomerContacts(result.Model.Properties.CustomerContacts)
					state.Name = pointer.ToString(result.Model.Name)
					state.Location = result.Model.Location
					state.Zones = result.Model.Zones
					state.ResourceGroupName = id.ResourceGroupName
					state.Tags = pointer.From(result.Model.Tags)
					state.ComputeCount = pointer.From(props.ComputeCount)
					state.DisplayName = props.DisplayName
					state.StorageCount = pointer.From(props.StorageCount)
					state.Shape = props.Shape
					state.MaintenanceWindow = FlattenMaintenanceWindow(props.MaintenanceWindow)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (ExadataInfraResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudExadataInfrastructures

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
	customerContacts := make([]cloudexadatainfrastructures.CustomerContact, 0, len(customerContactsList))
	for _, customerContact := range customerContactsList {
		customerContacts = append(customerContacts, cloudexadatainfrastructures.CustomerContact{
			Email: customerContact,
		})
	}
	return customerContacts
}

func ExpandDayOfWeekTo(daysOfWeek []string) []cloudexadatainfrastructures.DayOfWeek {
	daysOfWeekConverted := make([]cloudexadatainfrastructures.DayOfWeek, 0, len(daysOfWeek))
	for _, day := range daysOfWeek {
		daysOfWeekConverted = append(daysOfWeekConverted, cloudexadatainfrastructures.DayOfWeek{
			Name: cloudexadatainfrastructures.DayOfWeekName(day),
		})
	}
	return daysOfWeekConverted
}

func ExpandMonths(months []string) []cloudexadatainfrastructures.Month {
	monthsConverted := make([]cloudexadatainfrastructures.Month, 0, len(months))
	for _, month := range months {
		monthsConverted = append(monthsConverted, cloudexadatainfrastructures.Month{
			Name: cloudexadatainfrastructures.MonthName(month),
		})
	}
	return monthsConverted
}
