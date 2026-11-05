// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.WatcherClient

			var model WatcherModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionID := metadata.Client.Account.SubscriptionId
			accountID, _ := watcher.ParseAutomationAccountID(model.AutomationAccountID)
			id := watcher.NewWatcherID(subscriptionID, accountID.ResourceGroupName, accountID.AutomationAccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			tags := helpers.ExpandStringInterfaceMap(model.Tags)
			scriptParameters := helpers.ExpandStringInterfaceMap(model.ScriptParameters)

			param := watcher.Watcher{
				Properties: &watcher.WatcherProperties{
					Description:                 pointer.To(model.Description),
					ExecutionFrequencyInSeconds: pointer.To(model.ExecutionFrequencyInSeconds),
					ScriptName:                  pointer.To(model.ScriptName),
					ScriptParameters:            &scriptParameters,
					ScriptRunOn:                 pointer.To(model.ScriptRunOn),
				},
				Etag:     pointer.To(model.Etag),
				Location: pointer.To(model.Location),
				Tags:     &tags,
			}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
