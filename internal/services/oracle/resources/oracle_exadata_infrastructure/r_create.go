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

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := cloudexadatainfrastructures.CloudExadataInfrastructure{
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Zones:    zones.Expand(model.Zones),
				Properties: &cloudexadatainfrastructures.CloudExadataInfrastructureProperties{
					ComputeCount:     pointer.To(model.ComputeCount),
					DisplayName:      model.DisplayName,
					StorageCount:     pointer.To(model.StorageCount),
					Shape:            model.Shape,
					CustomerContacts: pointer.To(ExpandCustomerContacts(model.CustomerContacts)),
				},
			}
			if model.DatabaseServerType != "" {
				param.Properties.DatabaseServerType = pointer.To(model.DatabaseServerType)
			}
			if model.StorageServerType != "" {
				param.Properties.StorageServerType = pointer.To(model.StorageServerType)
			}
			if len(model.MaintenanceWindow) > 0 {
				param.Properties.MaintenanceWindow = &cloudexadatainfrastructures.MaintenanceWindow{
					DaysOfWeek:      pointer.To(ExpandDayOfWeekTo(model.MaintenanceWindow[0].DaysOfWeek)),
					HoursOfDay:      pointer.To(model.MaintenanceWindow[0].HoursOfDay),
					LeadTimeInWeeks: pointer.To(model.MaintenanceWindow[0].LeadTimeInWeeks),
					Months:          pointer.To(ExpandMonths(model.MaintenanceWindow[0].Months)),
					PatchingMode:    pointer.ToEnum[cloudexadatainfrastructures.PatchingMode](model.MaintenanceWindow[0].PatchingMode),
					Preference:      pointer.ToEnum[cloudexadatainfrastructures.Preference](model.MaintenanceWindow[0].Preference),
					WeeksOfMonth:    pointer.To(model.MaintenanceWindow[0].WeeksOfMonth),
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
