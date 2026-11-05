// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_powershell72_module

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.ModuleClientV2023

			var model AutomationPowerShell72ModuleModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			accountID, _ := module.ParseAutomationAccountID(model.AutomationAccountID)
			name := metadata.ResourceData.Get("name").(string)

			id := module.NewPowerShell72ModuleID(subscriptionId, accountID.ResourceGroupName, accountID.AutomationAccountName, name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.PowerShell72ModuleGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				// for existing global module do update instead of raising ImportAsExistsError
				isGlobal := existing.Model != nil && existing.Model.Properties != nil && pointer.From(existing.Model.Properties.IsGlobal)
				if !response.WasNotFound(existing.HttpResponse) && !isGlobal {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := module.ModuleCreateOrUpdateParameters{
				Properties: module.ModuleCreateOrUpdateProperties{
					ContentLink: expandPowerShell72ModuleLink(model.ModuleLink),
				},
				Tags: tags.Expand(model.Tags),
			}

			if _, err := client.PowerShell72ModuleCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			// the API returns 'done' but it's not actually finished provisioning yet
			// tracking issue: https://github.com/Azure/azure-rest-api-specs/pull/25435
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					string(module.ModuleProvisioningStateActivitiesStored),
					string(module.ModuleProvisioningStateConnectionTypeImported),
					string(module.ModuleProvisioningStateContentDownloaded),
					string(module.ModuleProvisioningStateContentRetrieved),
					string(module.ModuleProvisioningStateContentStored),
					string(module.ModuleProvisioningStateContentValidated),
					string(module.ModuleProvisioningStateCreated),
					string(module.ModuleProvisioningStateCreating),
					string(module.ModuleProvisioningStateModuleDataStored),
					string(module.ModuleProvisioningStateModuleImportRunbookComplete),
					string(module.ModuleProvisioningStateRunningImportModuleRunbook),
					string(module.ModuleProvisioningStateStartingImportModuleRunbook),
					string(module.ModuleProvisioningStateUpdating),
				},
				Target: []string{
					string(module.ModuleProvisioningStateSucceeded),
				},
				MinTimeout: 30 * time.Second,
				Refresh: func() (any, string, error) {
					resp, err2 := client.PowerShell72ModuleGet(ctx, id)
					if err2 != nil {
						return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err2)
					}

					provisioningState := "Unknown"
					if model := resp.Model; model != nil {
						if props := model.Properties; props != nil {
							if props.ProvisioningState != nil {
								provisioningState = string(*props.ProvisioningState)
							}
							if props.Error != nil && props.Error.Message != nil && *props.Error.Message != "" {
								return resp, provisioningState, errors.New(*props.Error.Message)
							}
							return resp, provisioningState, nil
						}
					}
					return resp, provisioningState, nil
				},
				Timeout: time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPowerShell72ModuleLink(m []ModuleLinkModel) module.ContentLink {
	moduleLink := m[0]
	if len(moduleLink.Hash) > 0 {
		hash := moduleLink.Hash[0]
		return module.ContentLink{
			Uri: &moduleLink.Uri,
			ContentHash: &module.ContentHash{
				Algorithm: hash.Algorithm,
				Value:     hash.Value,
			},
		}
	}

	return module.ContentLink{
		Uri: &moduleLink.Uri,
	}
}
