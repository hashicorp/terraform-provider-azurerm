// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_module

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/module"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationModuleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Module
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := module.NewModuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		// for existing global module do update instead of raising ImportAsExistsError
		isGlobal := existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.IsGlobal != nil && *existing.Model.Properties.IsGlobal
		if !response.WasNotFound(existing.HttpResponse) && !isGlobal {
			return tf.ImportAsExistsError("azurerm_automation_module", id.ID())
		}
	}

	parameters := module.ModuleCreateOrUpdateParameters{
		Properties: module.ModuleCreateOrUpdateProperties{
			ContentLink: expandModuleLink(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := waitForModuleProvisioningCompletion(ctx, client, id, d.Timeout(pluginsdk.TimeoutCreate)); err != nil {
		return err
	}

	return resourceAutomationModuleRead(d, meta)
}

func expandModuleLink(d *pluginsdk.ResourceData) module.ContentLink {
	inputs := d.Get("module_link").([]any)
	input := inputs[0].(map[string]any)
	uri := input["uri"].(string)

	hashes := input["hash"].([]any)

	if len(hashes) > 0 {
		hash := hashes[0].(map[string]any)
		return module.ContentLink{
			Uri: &uri,
			ContentHash: &module.ContentHash{
				Algorithm: hash["algorithm"].(string),
				Value:     hash["value"].(string),
			},
		}
	}

	return module.ContentLink{
		Uri: &uri,
	}
}
