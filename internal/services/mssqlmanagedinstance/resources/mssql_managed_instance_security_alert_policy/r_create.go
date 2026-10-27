// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_security_alert_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedserversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlManagedInstanceSecurityAlertPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	alertPolicy := expandManagedServerSecurityAlertPolicy(d)

	managedInstanceId := commonids.NewSqlManagedInstanceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("managed_instance_name").(string))

	err := client.CreateOrUpdateThenPoll(ctx, managedInstanceId, *alertPolicy)
	if err != nil {
		return fmt.Errorf("updating managed instance security alert policy: %v", err)
	}

	result, err := client.Get(ctx, managedInstanceId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", managedInstanceId, err)
	}

	if result.Model == nil || result.Model.Name == nil {
		return fmt.Errorf("reading %s", managedInstanceId)
	}

	id := parse.NewManagedInstancesSecurityAlertPolicyID(subscriptionId, managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, *result.Model.Name)

	d.SetId(id.ID())

	return resourceMsSqlManagedInstanceSecurityAlertPolicyRead(d, meta)
}

func expandManagedServerSecurityAlertPolicy(d *pluginsdk.ResourceData) *managedserversecurityalertpolicies.ManagedServerSecurityAlertPolicy {
	state := managedserversecurityalertpolicies.SecurityAlertsPolicyStateDisabled

	if d.Get("enabled").(bool) {
		state = managedserversecurityalertpolicies.SecurityAlertsPolicyStateEnabled
	}

	policy := managedserversecurityalertpolicies.ManagedServerSecurityAlertPolicy{
		Properties: &managedserversecurityalertpolicies.SecurityAlertsPolicyProperties{
			State: state,
		},
	}

	props := policy.Properties

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
		props.RetentionDays = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		props.StorageAccountAccessKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("storage_endpoint"); ok {
		props.StorageEndpoint = pointer.To(v.(string))
	}

	return &policy
}
