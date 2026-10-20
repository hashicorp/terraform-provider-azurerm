// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApplicationInsightsWorkbookModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppInsights.WorkbookClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := workbooksapis.NewWorkbookID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.WorkbooksGet(ctx, id, workbooksapis.WorkbooksGetOperationOptions{CanFetchContent: pointer.To(true)})
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			properties := &workbooksapis.Workbook{
				Identity: identityValue,
				Kind:     pointer.To(workbooksapis.WorkbookSharedTypeKindShared),
				Location: location.Normalize(model.Location),
				Properties: &workbooksapis.WorkbookProperties{
					Category:       model.Category,
					DisplayName:    model.DisplayName,
					SerializedData: model.DataJson,
					SourceId:       &model.SourceId,
				},

				Tags: &model.Tags,
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.StorageContainerId != "" {
				properties.Properties.StorageUri = &model.StorageContainerId
			}

			if _, err := client.WorkbooksCreateOrUpdate(ctx, id, *properties, workbooksapis.WorkbooksCreateOrUpdateOperationOptions{SourceId: &model.SourceId}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}
