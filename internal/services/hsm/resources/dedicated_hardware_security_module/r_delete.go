// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_hardware_security_module

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	hsmClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/hsm/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHardwareSecurityModuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := hsmClient.From(meta).DedicatedHsmClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := dedicatedhsms.ParseDedicatedHSMID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DedicatedHsmDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
