// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_hardware_security_module

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHardwareSecurityModuleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhsms.ParseDedicatedHSMID(d.Id())
	if err != nil {
		return err
	}

	parameters := dedicatedhsms.DedicatedHsmPatchParameters{}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.DedicatedHsmUpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDedicatedHardwareSecurityModuleRead(d, meta)
}
