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
				RequiredWith: []string{"storage_endpoint"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				RequiredWith: []string{"storage_account_access_key"},
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

	log.Printf("[INFO] preparing arguments for mssql server security alert policy creation")

	payload := serversecurityalertpolicies.ServerSecurityAlertPolicy{}
	resourceGroupName := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	state := d.Get("state").(string)

	props := &serversecurityalertpolicies.SecurityAlertsPolicyProperties{
		State: serversecurityalertpolicies.SecurityAlertsPolicyState(state),
	}

	var disabledAlerts *[]string
	var emailAddresses *[]string
	var emailAdmins *bool
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

	// NOTE: The API defaults to 'true' for the 'EmailAccountAdmins'
	// property, the provider defaults to 'false'...
	if v, ok := d.GetOk("email_account_admins"); ok {
		emailAdmins = pointer.To(v.(bool))
	}
	props.EmailAccountAdmins = emailAdmins

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

	err := client.CreateOrUpdateThenPoll(ctx, serverId, payload)
	if err != nil {
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

	result, err := client.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			log.Printf("[WARN] mssql server security alert policy %s: not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request to mssql server security alert policy: %+v", err)
	}

	model := result.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: Properties was nil", id)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)
	d.Set("state", string(props.State))

	disabledAlerts := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
	if props.DisabledAlerts != nil {
		for _, v := range *props.DisabledAlerts {
			if v != "" {
				disabledAlerts.Add(v)
			}
		}
	}
	d.Set("disabled_alerts", disabledAlerts)

	var emailAdmins bool
	if props.EmailAccountAdmins != nil {
		emailAdmins = *props.EmailAccountAdmins
	}
	d.Set("email_account_admins", emailAdmins)

	emailAddresses := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
	if props.EmailAddresses != nil {
		for _, v := range *props.EmailAddresses {
			if v != "" {
				emailAddresses.Add(v)
			}
		}
	}
	d.Set("email_addresses", emailAddresses)

	var retentionDays int
	if props.RetentionDays != nil {
		retentionDays = int(*props.RetentionDays)
	}
	d.Set("retention_days", retentionDays)

	var storageEndpoint string
	if props.StorageEndpoint != nil {
		storageEndpoint = *props.StorageEndpoint
	}
	d.Set("storage_endpoint", storageEndpoint)

	// NOTE: 'storage_account_access_key' field is not returned by the API
	// so we need to pull it from the state...
	var accessKey string
	if v, ok := d.GetOk("storage_account_access_key"); ok {
		accessKey = v.(string)
	}
	d.Set("storage_account_access_key", accessKey)

	return nil
}

func resourceMsSqlServerSecurityAlertPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for mssql server security alert policy update")

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	policy, err := client.Get(ctx, serverId)
	if err != nil {
		return fmt.Errorf("retrieving mssql server security alert policy %s: %+v", id, err)
	}

	model := policy.Model
	if model == nil {
		return fmt.Errorf("retrieving mssql server security alert policy %s: model was nil", id)
	}

	if model.Name == nil {
		return fmt.Errorf("reading mssql server security alert policy %s: name was nil", id)
	}

	if model.Properties == nil {
		return fmt.Errorf("reading mssql server security alert policy %s: properties was nil", id)
	}

	payload := serversecurityalertpolicies.ServerSecurityAlertPolicy{}
	props := model.Properties

	if d.HasChange("state") {
		props.State = serversecurityalertpolicies.SecurityAlertsPolicyState(d.Get("state").(string))
	}

	if d.HasChange("disabled_alerts") {
		disabledAlerts := make([]string, 0)
		if v, ok := d.GetOk("disabled_alerts"); ok {
			for _, v := range v.(*pluginsdk.Set).List() {
				disabledAlerts = append(disabledAlerts, v.(string))
			}
		}
		props.DisabledAlerts = pointer.To(disabledAlerts)
	}

	if d.HasChange("email_addresses") {
		emailAddresses := make([]string, 0)
		if v, ok := d.GetOk("email_addresses"); ok {
			for _, v := range v.(*pluginsdk.Set).List() {
				emailAddresses = append(emailAddresses, v.(string))
			}
		}
		props.EmailAddresses = pointer.To(emailAddresses)
	}

	if d.HasChange("email_account_admins") {
		var emailAdmins *bool
		if v, ok := d.GetOk("email_account_admins"); ok {
			emailAdmins = pointer.To(v.(bool))
		}
		props.EmailAccountAdmins = emailAdmins
	}

	if d.HasChange("retention_days") {
		var retentionDays *int64
		if v, ok := d.GetOk("retention_days"); ok {
			retentionDays = pointer.To(int64(v.(int)))
		}
		props.RetentionDays = retentionDays
	}

	if d.HasChange("storage_endpoint") {
		props.StorageEndpoint = nil
		if v, ok := d.GetOk("storage_endpoint"); ok {
			props.StorageEndpoint = pointer.To(v.(string))
		}
	}

	// NOTE: 'storage_account_access_key' it is not returned
	// by the API, so we need to get it from state...
	props.StorageAccountAccessKey = nil
	if v, ok := d.GetOk("storage_account_access_key"); ok {
		props.StorageAccountAccessKey = pointer.To(v.(string))
	}

	payload.Properties = props

	err = client.CreateOrUpdateThenPoll(ctx, serverId, payload)
	if err != nil {
		return fmt.Errorf("updating mssql server security alert policy: %+v", err)
	}

	return resourceMsSqlServerSecurityAlertPolicyRead(d, meta)
}

func resourceMsSqlServerSecurityAlertPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] deleting mssql server security alert policy")

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
		return fmt.Errorf("updating mssql server security alert policy: %+v", err)
	}

	if _, err = client.Get(ctx, serverId); err != nil {
		return fmt.Errorf("deleting mssql server security alert policy: %+v", err)
	}

	return nil
}
