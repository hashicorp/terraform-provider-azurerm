// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO 4.0 - consider/investigate inlining this within the mssql_server resource now that it exists.

func resourceMsSqlServerSecurityAlertPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlServerSecurityAlertPolicyCreateUpdate,
		Read:   resourceMsSqlServerSecurityAlertPolicyRead,
		Update: resourceMsSqlServerSecurityAlertPolicyCreateUpdate,
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
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.SecurityAlertPolicyStateDisabled),
					string(sql.SecurityAlertPolicyStateEnabled),
					string(sql.SecurityAlertPolicyStateNew),
				}, false),
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

func resourceMsSqlServerSecurityAlertPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for mssql server security alert policy creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	alertPolicy := expandSecurityAlertPolicy(d)

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, serverName, *alertPolicy)
	if err != nil {
		return fmt.Errorf("updataing mssql server security alert policy: %v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of mssql server security alert policy (server %q, resource group %q): %+v", serverName, resourceGroupName, err)
	}

	result, err := client.Get(ctx, resourceGroupName, serverName)
	if err != nil {
		return fmt.Errorf("retrieving mssql server security alert policy (server %q, resource group %q): %+v", serverName, resourceGroupName, err)
	}
	if result.Name == nil {
		return fmt.Errorf("reading mssql server security alert policy name (server %q, resource group %q)", serverName, resourceGroupName)
	}
	id := parse.NewServerSecurityAlertPolicyID(subscriptionId, resourceGroupName, serverName, *result.Name)

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

	result, err := client.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[WARN] mssql server security alert policy %v not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request to mssql server security alert policy: %+v", err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	if props := result.SecurityAlertsPolicyProperties; props != nil {
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

func resourceMsSqlServerSecurityAlertPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] deleting mssql server security alert policy.")

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	disabledPolicy := sql.ServerSecurityAlertPolicy{
		SecurityAlertsPolicyProperties: &sql.SecurityAlertsPolicyProperties{
			State: sql.SecurityAlertsPolicyStateDisabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, disabledPolicy)
	if err != nil {
		return fmt.Errorf("updataing mssql server security alert policy: %v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of mssql server security alert policy (server %q, resource group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if _, err = client.Get(ctx, id.ResourceGroup, id.ServerName); err != nil {
		return fmt.Errorf("deleting mssql server security alert policy: %v", err)
	}

	return nil
}

func expandSecurityAlertPolicy(d *pluginsdk.ResourceData) *sql.ServerSecurityAlertPolicy {
	state := sql.SecurityAlertsPolicyState(d.Get("state").(string))

	policy := sql.ServerSecurityAlertPolicy{
		SecurityAlertsPolicyProperties: &sql.SecurityAlertsPolicyProperties{
			State: state,
		},
	}

	props := policy.SecurityAlertsPolicyProperties

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
		props.EmailAccountAdmins = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("retention_days"); ok {
		props.RetentionDays = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		props.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_endpoint"); ok {
		props.StorageEndpoint = utils.String(v.(string))
	}

	return &policy
}
