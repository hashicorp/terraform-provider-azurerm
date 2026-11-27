// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_firewall_rule

import (
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/ipfirewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseFirewallRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).FirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := ipfirewallrules.NewFirewallRuleID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName, id.FirewallRuleName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_synapse_firewall_rule", id.ID())
			}
		}
	}

	parameters := synapse.IPFirewallRuleInfo{
		IPFirewallRuleProperties: &synapse.IPFirewallRuleProperties{
			StartIPAddress: pointer.To(d.Get("start_ip_address").(string)),
			EndIPAddress:   pointer.To(d.Get("end_ip_address").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.WorkspaceName, id.FirewallRuleName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation/update of %s: %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("internal-error: context had no deadline")
	}

	// The firewall is not taking effect immediately after firewall creation.
	// Firewall has a cache and will refresh every 1 minute, so if requests sent before firewall refreshes, it will meet ClientIpAddressNotAuthorized.
	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/21516
	stateChangeConf := &pluginsdk.StateChangeConf{
		Pending: []string{string(synapse.ProvisioningStateProvisioning)},
		Target:  []string{string(synapse.ProvisioningStateSucceeded)},
		Refresh: func() (result any, state string, err error) {
			resp, err := client.Get(ctx, id.ResourceGroupName, id.WorkspaceName, id.FirewallRuleName)
			if err != nil {
				return nil, "Error", err
			}
			return resp, string(resp.ProvisioningState), err
		},
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(deadline),
	}

	if _, err = stateChangeConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be ready", id)
	}

	return resourceSynapseFirewallRuleRead(d, meta)
}
