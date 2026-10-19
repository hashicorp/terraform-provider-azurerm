// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.AutomationAccount
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	identityVal, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	enc, err := expandEncryption(d.Get("encryption").([]any))
	if err != nil {
		return fmt.Errorf("expanding `encryption`: %v", err)
	}

	parameters := automationaccount.AutomationAccountUpdateParameters{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Identity: identityVal,
		Properties: &automationaccount.AutomationAccountUpdateProperties{
			Sku: &automationaccount.Sku{
				Name: automationaccount.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: pointer.To(d.Get("public_network_access_enabled").(bool)),
			Encryption:          enc,
		},
	}

	if d.HasChange("local_authentication_enabled") {
		parameters.Properties.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))
	}

	if tagsVal := tags.Expand(d.Get("tags").(map[string]any)); tagsVal != nil {
		parameters.Tags = tagsVal
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationAccountRead(d, meta)
}
