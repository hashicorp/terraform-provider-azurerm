// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package advanced_threat_protection

import (
	"fmt"
	"log"

	//nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAdvancedThreatProtectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AdvancedThreatProtectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.TargetResourceID)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("Advanced Threat Protection was not found for %q: %+v", id.TargetResourceID, err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Advanced Threat Protection status for %q: %+v", id.TargetResourceID, err)
	}

	d.Set("target_resource_id", id.TargetResourceID)
	if atpp := resp.AdvancedThreatProtectionProperties; atpp != nil {
		d.Set("enabled", pointer.From(resp.IsEnabled))
	}

	return nil
}
