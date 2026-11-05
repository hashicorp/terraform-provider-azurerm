// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_module

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationModuleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Module
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := module.ParseModuleID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	// Start from the existing content link so that fields managed via
	// ignore_changes retain their server-side value.
	var contentLink *module.ContentLink
	if existing.Model != nil && existing.Model.Properties != nil {
		contentLink = existing.Model.Properties.ContentLink
	}

	if d.HasChange("module_link") {
		expanded := expandModuleLink(d)
		contentLink = &expanded
	}

	parameters := module.ModuleUpdateParameters{
		Name: &id.ModuleName,
		Properties: &module.ModuleUpdateProperties{
			ContentLink: contentLink,
		},
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := waitForModuleProvisioningCompletion(ctx, client, *id, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceAutomationModuleRead(d, meta)
}
