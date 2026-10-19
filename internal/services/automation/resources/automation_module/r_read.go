// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_module

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationModuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Module
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := module.ParseModuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ModuleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	return nil
}

// waitForModuleProvisioningCompletion polls the module until it reaches the Succeeded
// provisioning state. The API returns 'done' before provisioning is actually complete.
// Tracking issue: https://github.com/Azure/azure-rest-api-specs/pull/25435
func waitForModuleProvisioningCompletion(ctx context.Context, client *module.ModuleClient, id module.ModuleId, timeout time.Duration) error {
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
		Timeout:    timeout,
		Refresh: func() (any, string, error) {
			resp, err := client.Get(ctx, id)
			if err != nil {
				return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
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
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)
	}

	return nil
}
