// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package advanced_threat_protection

import (
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAdvancedThreatProtectionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAdvancedThreatProtectionId(d.Get("target_resource_id").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			server, err := client.Get(ctx, id.TargetResourceID)
			if err != nil {
				if !response.WasNotFound(server.Response.Response) {
					return fmt.Errorf("checking for presence of existing Advanced Threat Protection for %q: %+v", id.TargetResourceID, err)
				}
			}

			if server.ID != nil && *server.ID != "" && server.IsEnabled != nil && *server.IsEnabled {
				return tf.ImportAsExistsError("azurerm_advanced_threat_protection", id.ID())
			}
		}
	}

	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: pointer.To(d.Get("enabled").(bool)),
		},
	}

	if _, err := client.Create(ctx, id.TargetResourceID, setting); err != nil {
		return fmt.Errorf("updating Advanced Threat protection for %q: %+v", id.TargetResourceID, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}

	// the API appears to be eventually consistent, tracked on https://github.com/Azure/azure-rest-api-specs/issues/25232
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"diff"},
		Target:  []string{"consistent"},
		Refresh: func() (result any, state string, err error) {
			resp, err := client.Get(ctx, id.TargetResourceID)
			if err != nil {
				return resp, "error", err
			}
			if atpp := resp.AdvancedThreatProtectionProperties; atpp != nil {
				respEnabled := atpp.IsEnabled != nil && *atpp.IsEnabled
				if respEnabled == d.Get("enabled").(bool) {
					return resp, "consistent", nil
				} else {
					return resp, "diff", nil
				}
			}
			return resp, "error", errors.New("properties was nil")
		},
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of advanced threat protection: %+v", err)
	}

	d.SetId(id.ID())
	return resourceAdvancedThreatProtectionRead(d, meta)
}
