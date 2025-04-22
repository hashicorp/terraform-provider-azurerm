// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// TODO 4.0 - consider/investigate inlining this within the mssql_server resource now that it exists.

func resourceMsSqlServerSecurityAlertPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlServerSecurityAlertPolicyCreate,
		Read:   resourceMsSqlServerSecurityAlertPolicyRead,
		Update: resourceMsSqlServerSecurityAlertPolicyUpdate,
		Delete: resourceMsSqlServerSecurityAlertPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServerSecurityAlertPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"disabled_alerts": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Set:      pluginsdk.HashString,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"Sql_Injection",
						"Sql_Injection_Vulnerability",
						"Access_Anomaly",
						"Data_Exfiltration",
						"Unsafe_Action",
					}, false),
				},
			},

			"email_account_admins": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"email_addresses": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				Set: pluginsdk.HashString,
			},

			"retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"state": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(serversecurityalertpolicies.PossibleValuesForSecurityAlertsPolicyState(),
					false),
			},

			"storage_account_access_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceMsSqlServerSecurityAlertPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for mssql server security alert policy creation.")

	alertPolicy := serversecurityalertpolicies.ServerSecurityAlertPolicy{}
	resourceGroupName := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	props := &serversecurityalertpolicies.SecurityAlertsPolicyProperties{
		State: serversecurityalertpolicies.SecurityAlertsPolicyState(d.Get("state").(string)),
	}

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

	if v, ok := d.GetOk("email_account_admins"); ok {
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

	alertPolicy.Properties = props
	serverId := commonids.NewSqlServerID(subscriptionId, resourceGroupName, serverName)

	err := client.CreateOrUpdateThenPoll(ctx, serverId, alertPolicy)
	if err != nil {
		return fmt.Errorf("creating mssql server security alert policy: %v", err)
	}

	policy, err := client.Get(ctx, serverId)
	if err != nil {
		return fmt.Errorf("retrieving mssql server security alert policy (server %q, resource group %q): %+v", serverName, resourceGroupName, err)
	}

	model := policy.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", serverId)
	}

	if model.Name == nil {
		return fmt.Errorf("reading mssql server security alert policy name (server %q, resource group %q)", serverName, resourceGroupName)
	}

	id := parse.NewServerSecurityAlertPolicyID(subscriptionId, resourceGroupName, serverName, *model.Name)

	d.SetId(id.ID())

	return resourceMsSqlServerSecurityAlertPolicyRead(d, meta)
}

func resourceMsSqlServerSecurityAlertPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] reading mssql server security alert policy")

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	policy, err := client.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(policy.HttpResponse) {
			log.Printf("[WARN] mssql server security alert policy %v not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request to mssql server security alert policy: %+v", err)
	}

	model := policy.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	if props := model.Properties; props != nil {
		d.Set("state", string(props.State))

		if props.DisabledAlerts != nil {
			disabledAlerts := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
			for _, v := range *props.DisabledAlerts {
				if v != "" {
					disabledAlerts.Add(v)
				}
			}

			d.Set("disabled_alerts", disabledAlerts)
		}

		if props.EmailAccountAdmins != nil {
			d.Set("email_account_admins", props.EmailAccountAdmins)
		}

		if props.EmailAddresses != nil {
			emailAddresses := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
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

func resourceMsSqlServerSecurityAlertPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for mssql server security alert policy update.")

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	policy, err := client.Get(ctx, serverId)
	if err != nil {
		return fmt.Errorf("retrieving mssql server security alert policy (server %q, resource group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	model := policy.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", serverId)
	}

	if model.Name == nil {
		return fmt.Errorf("reading mssql server security alert policy name (server %q, resource group %q)", id.ServerName, id.ResourceGroup)
	}

	if model.Properties == nil {
		return fmt.Errorf("reading mssql server security alert policy properties (server %q, resource group %q)", id.ServerName, id.ResourceGroup)
	}

	alertPolicy := serversecurityalertpolicies.ServerSecurityAlertPolicy{}
	props := model.Properties

	if d.HasChange("state") {
		props.State = serversecurityalertpolicies.SecurityAlertsPolicyState(d.Get("state").(string))
	}

	if d.HasChange("disabled_alerts") {
		if v, ok := d.GetOk("disabled_alerts"); ok {
			disabledAlerts := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				disabledAlerts = append(disabledAlerts, v.(string))
			}
			props.DisabledAlerts = &disabledAlerts
		}
	}

	if d.HasChange("email_addresses") {
		if v, ok := d.GetOk("email_addresses"); ok {
			emailAddresses := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				emailAddresses = append(emailAddresses, v.(string))
			}
			props.EmailAddresses = &emailAddresses
		}
	}

	if d.HasChange("email_account_admins") {
		if v, ok := d.GetOk("email_account_admins"); ok {
			props.EmailAccountAdmins = pointer.To(v.(bool))
		}
	}

	if d.HasChange("retention_days") {
		if v, ok := d.GetOk("retention_days"); ok {
			props.RetentionDays = pointer.To(int64(v.(int)))
		}
	}

	if d.HasChange("storage_account_access_key") {
		if v, ok := d.GetOk("storage_account_access_key"); ok {
			props.StorageAccountAccessKey = pointer.To(v.(string))
		}
	}

	if d.HasChange("storage_endpoint") {
		if v, ok := d.GetOk("storage_endpoint"); ok {
			props.StorageEndpoint = pointer.To(v.(string))
		}
	}

	alertPolicy.Properties = props

	err = client.CreateOrUpdateThenPoll(ctx, serverId, alertPolicy)
	if err != nil {
		return fmt.Errorf("updating mssql server security alert policy: %v", err)
	}

	return resourceMsSqlServerSecurityAlertPolicyRead(d, meta)
}

func resourceMsSqlServerSecurityAlertPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] deleting mssql server security alert policy.")

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	disabledPolicy := serversecurityalertpolicies.ServerSecurityAlertPolicy{
		Properties: &serversecurityalertpolicies.SecurityAlertsPolicyProperties{
			State: serversecurityalertpolicies.SecurityAlertsPolicyStateDisabled,
		},
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	err = client.CreateOrUpdateThenPoll(ctx, serverId, disabledPolicy)
	if err != nil {
		return fmt.Errorf("updating mssql server security alert policy: %v", err)
	}

	if _, err = client.Get(ctx, serverId); err != nil {
		return fmt.Errorf("deleting mssql server security alert policy: %v", err)
	}

	return nil
}
