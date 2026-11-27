// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package advanced_threat_protection

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAdvancedThreatProtectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.AdvancedThreatProtectionID(d.Id())
	if err != nil {
		return err
	}

	// there is no delete.. so lets just do best effort and set it to false?
	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: pointer.To(false),
		},
	}

	if _, err := client.Create(ctx, id.TargetResourceID, setting); err != nil {
		return fmt.Errorf("removing Advanced Threat Protection for %q: %+v", id.TargetResourceID, err)
	}

	return nil
}
