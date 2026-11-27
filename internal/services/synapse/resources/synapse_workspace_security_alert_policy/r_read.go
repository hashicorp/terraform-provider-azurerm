// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_security_alert_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceSecurityAlertPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceSecurityAlertPolicyClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.WorkspaceSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] synapse %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("synapse_workspace_id", workspaceId.ID())

	if props := resp.ServerSecurityAlertPolicyProperties; props != nil {
		d.Set("policy_state", string(props.State))

		if props.DisabledAlerts != nil {
			disabledAlerts := pluginsdk.NewSet(pluginsdk.HashString, []any{})
			for _, v := range *props.DisabledAlerts {
				if v != "" {
					disabledAlerts.Add(v)
				}
			}

			d.Set("disabled_alerts", disabledAlerts)
		}

		if props.EmailAccountAdmins != nil {
			d.Set("email_account_admins_enabled", props.EmailAccountAdmins)
		}

		if props.EmailAddresses != nil {
			emailAddresses := pluginsdk.NewSet(pluginsdk.HashString, []any{})
			for _, v := range *props.EmailAddresses {
				if v != "" {
					emailAddresses.Add(v)
				}
			}

			d.Set("email_addresses", emailAddresses)
		}

		if props.RetentionDays != nil {
			d.Set("retention_days", int(*props.RetentionDays))
		}

		if v, ok := d.GetOk("storage_account_access_key"); ok {
			d.Set("storage_account_access_key", v)
		}

		if props.StorageEndpoint != nil {
			d.Set("storage_endpoint", props.StorageEndpoint)
		}
	}

	return nil
}
