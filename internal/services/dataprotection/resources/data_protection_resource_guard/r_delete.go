// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_resource_guard

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionResourceGuardDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.ResourceGuardClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resourceguardresources.ParseResourceGuardID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.ResourceGuardsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
