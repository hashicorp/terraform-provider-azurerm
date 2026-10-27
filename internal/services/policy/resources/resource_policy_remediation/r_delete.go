// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_policy_remediation

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourcePolicyRemediationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseScopedRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Policy Scoped Remediation ID: %+v", err)
	}

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.GetAtResource(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := helpers.WaitForRemediationToDelete(
		ctx,
		existing.Model.Properties,
		id.ID(),
		d.Timeout(pluginsdk.TimeoutDelete),
		func() error {
			_, err := client.CancelAtResource(ctx, *id)
			return err
		},
		resourcePolicyRemediationCancellationRefreshFunc(ctx, client, *id),
	); err != nil {
		return err
	}

	_, err = client.DeleteAtResource(ctx, *id)

	return err
}
