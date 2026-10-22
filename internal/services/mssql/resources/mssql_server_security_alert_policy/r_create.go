// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_security_alert_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerSecurityAlertPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	payload := serversecurityalertpolicies.ServerSecurityAlertPolicy{}
	resourceGroupName := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	state := d.Get("state").(string)

	props := &serversecurityalertpolicies.SecurityAlertsPolicyProperties{
		State: serversecurityalertpolicies.SecurityAlertsPolicyState(state),
	}

	var disabledAlerts *[]string
	var emailAddresses *[]string
	var retentionDays *int64
	var storageAccountAccessKey *string
	var storageEndpoint *string

	if v, ok := d.GetOk("disabled_alerts"); ok {
		disabled := make([]string, 0)
		for _, v := range v.(*pluginsdk.Set).List() {
			disabled = append(disabled, v.(string))
		}
		disabledAlerts = pointer.To(disabled)
	}
	props.DisabledAlerts = disabledAlerts

	if v, ok := d.GetOk("email_addresses"); ok {
		emails := make([]string, 0)
		for _, v := range v.(*pluginsdk.Set).List() {
			emails = append(emails, v.(string))
		}
		emailAddresses = pointer.To(emails)
	}
	props.EmailAddresses = emailAddresses

	props.EmailAccountAdmins = pointer.To(d.Get("email_account_admins_enabled").(bool))

	if v, ok := d.GetOk("retention_days"); ok {
		retentionDays = pointer.To(int64(v.(int)))
	}
	props.RetentionDays = retentionDays

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		storageAccountAccessKey = pointer.To(v.(string))
	}
	props.StorageAccountAccessKey = storageAccountAccessKey

	if v, ok := d.GetOk("storage_endpoint"); ok {
		storageEndpoint = pointer.To(v.(string))
	}
	props.StorageEndpoint = storageEndpoint

	payload.Properties = props
	serverId := commonids.NewSqlServerID(subscriptionId, resourceGroupName, serverName)

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, serverId, payload); err != nil {
		return fmt.Errorf("creating mssql server security alert policy: %+v", err)
	}

	result, err := client.Get(ctx, serverId)
	if err != nil {
		return fmt.Errorf("retrieving mssql server security alert policy %s: %+v", serverId, err)
	}

	model := result.Model
	if model == nil {
		return fmt.Errorf("retrieving mssql server security alert policy %s: model was nil", serverId)
	}

	if model.Name == nil {
		return fmt.Errorf("retrieving mssql server security alert policy %s: name was nil", serverId)
	}

	id := parse.NewServerSecurityAlertPolicyID(subscriptionId, resourceGroupName, serverName, *model.Name)

	d.SetId(id.ID())

	return resourceMsSqlServerSecurityAlertPolicyRead(d, meta)
}
