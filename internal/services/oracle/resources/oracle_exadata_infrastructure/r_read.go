// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exadata_infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
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
				state.Zones = zones.Flatten(&model.Zones)
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.CustomerContacts = FlattenCustomerContacts(result.Model.Properties.CustomerContacts)
					state.Name = pointer.From(result.Model.Name)
					state.ResourceGroupName = id.ResourceGroupName
					state.ComputeCount = pointer.From(props.ComputeCount)
					state.DisplayName = props.DisplayName
					state.StorageCount = pointer.From(props.StorageCount)
					state.Shape = props.Shape
					state.MaintenanceWindow = FlattenMaintenanceWindow(props.MaintenanceWindow)
					state.DatabaseServerType = pointer.From(props.DatabaseServerType)
					state.StorageServerType = pointer.From(props.StorageServerType)
				}
			}

			return metadata.Encode(&state)
		},
	}
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
