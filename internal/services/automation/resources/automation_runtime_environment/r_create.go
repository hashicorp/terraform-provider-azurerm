// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.RuntimeEnvironment

			var model AutomationRuntimeEnvironmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionID := metadata.Client.Account.SubscriptionId

			autAccId, err := runtimeenvironment.ParseAutomationAccountIDInsensitively(model.AutomationAccountId)
			if err != nil {
				return err
			}

			id := runtimeenvironment.NewRuntimeEnvironmentID(subscriptionID, autAccId.ResourceGroupName, autAccId.AutomationAccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := runtimeenvironment.RuntimeEnvironment{
				Location: location.Normalize(model.Location),
				Name:     pointer.To(model.Name),
				Properties: &runtimeenvironment.RuntimeEnvironmentProperties{
					Runtime: &runtimeenvironment.RuntimeProperties{
						Language: pointer.To(model.RuntimeLanguage),
						Version:  pointer.To(model.RuntimeVersion),
					},
				},
			}

			if model.Tags != nil {
				req.Tags = &model.Tags
			}

			if model.RuntimeDefaultPackages != nil {
				req.Properties.DefaultPackages = pointer.To(model.RuntimeDefaultPackages)
			}

			if model.Description != "" {
				req.Properties.Description = pointer.To(model.Description)
			}

			if _, err = client.Create(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
