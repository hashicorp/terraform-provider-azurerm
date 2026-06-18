// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/ipfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseFirewallRuleCreateUpdate,
		Read:   resourceSynapseFirewallRuleRead,
		Update: resourceSynapseFirewallRuleCreateUpdate,
		Delete: resourceSynapseFirewallRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallRuleName,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"start_ip_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"end_ip_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourceSynapseFirewallRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.FirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := ipfirewallrules.NewFirewallRuleID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_synapse_firewall_rule", id.ID())
			}
		}
	}

	parameters := ipfirewallrules.IPFirewallRuleInfo{
		Properties: &ipfirewallrules.IPFirewallRuleProperties{
			StartIPAddress: pointer.To(d.Get("start_ip_address").(string)),
			EndIPAddress:   pointer.To(d.Get("end_ip_address").(string)),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}

	// The firewall is not taking effect immediately after firewall creation.
	// Firewall has a cache and will refresh every 1 minute, so if requests sent before firewall refreshes, it will meet ClientIpAddressNotAuthorized.
	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/21516
	stateChangeConf := &pluginsdk.StateChangeConf{
		Pending: []string{string(ipfirewallrules.ProvisioningStateProvisioning)},
		Target:  []string{string(ipfirewallrules.ProvisioningStateSucceeded)},
		Refresh: func() (result interface{}, state string, err error) {
			resp, err := client.Get(ctx, id)
			if err != nil {
				return nil, "Error", err
			}
			provisioningState := ""
			if model := resp.Model; model != nil && model.Properties != nil && model.Properties.ProvisioningState != nil {
				provisioningState = string(*model.Properties.ProvisioningState)
			}
			return resp, provisioningState, nil
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

func resourceSynapseFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ipfirewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading Synapse Firewall Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID()
	d.Set("name", id.FirewallRuleName)
	d.Set("synapse_workspace_id", workspaceId)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("start_ip_address", props.StartIPAddress)
			d.Set("end_ip_address", props.EndIPAddress)
		}
	}

	return nil
}

func resourceSynapseFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ipfirewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
