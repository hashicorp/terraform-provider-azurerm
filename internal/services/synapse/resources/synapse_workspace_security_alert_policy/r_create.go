// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_security_alert_policy

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceSecurityAlertPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceSecurityAlertPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	// TODO: import check?

	id := parse.NewWorkspaceSecurityAlertPolicyID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, "Default")

	alertPolicy := expandServerSecurityAlertPolicy(d)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, *alertPolicy)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	return resourceSynapseWorkspaceSecurityAlertPolicyRead(d, meta)
}

func expandServerSecurityAlertPolicy(d *pluginsdk.ResourceData) *synapse.ServerSecurityAlertPolicy {
	policyState := synapse.SecurityAlertPolicyState(d.Get("policy_state").(string))

	policy := synapse.ServerSecurityAlertPolicy{
		ServerSecurityAlertPolicyProperties: &synapse.ServerSecurityAlertPolicyProperties{
			State: policyState,
		},
	}

	props := policy.ServerSecurityAlertPolicyProperties

	if v, ok := d.GetOk("disabled_alerts"); ok {
		disabledAlerts := make([]string, 0)
		for _, v := range v.(*pluginsdk.Set).List() {
			disabledAlerts = append(disabledAlerts, v.(string))
		}
		props.DisabledAlerts = &disabledAlerts
	}

	if v, ok := d.GetOk("email_addresses"); ok {
		emailAddresses := make([]string, 0)
		for _, v := range v.(*pluginsdk.Set).List() {
			emailAddresses = append(emailAddresses, v.(string))
		}
		props.EmailAddresses = &emailAddresses
	}

	if v, ok := d.GetOk("email_account_admins_enabled"); ok {
		props.EmailAccountAdmins = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("retention_days"); ok {
		props.RetentionDays = pointer.To(int32(v.(int)))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		props.StorageAccountAccessKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("storage_endpoint"); ok {
		props.StorageEndpoint = pointer.To(v.(string))
	}

	return &policy
}
