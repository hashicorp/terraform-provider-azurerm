// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_security_alert_policy

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSqlPoolSecurityAlertPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SqlPoolSecurityAlertPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	sqlPoolId, err := parse.SqlPoolID(d.Get("sql_pool_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSqlPoolSecurityAlertPolicyID(sqlPoolId.SubscriptionId, sqlPoolId.ResourceGroup, sqlPoolId.WorkspaceName, sqlPoolId.Name, "Default")

	// TODO: import check?

	alertPolicy := expandSQLPoolSecurityAlertPolicy(d)

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, *alertPolicy); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSynapseSqlPoolSecurityAlertPolicyRead(d, meta)
}

func expandSQLPoolSecurityAlertPolicy(d *pluginsdk.ResourceData) *synapse.SQLPoolSecurityAlertPolicy {
	policyState := synapse.SecurityAlertPolicyState(d.Get("policy_state").(string))

	policy := synapse.SQLPoolSecurityAlertPolicy{
		SecurityAlertPolicyProperties: &synapse.SecurityAlertPolicyProperties{
			State: policyState,
		},
	}

	props := policy.SecurityAlertPolicyProperties

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
