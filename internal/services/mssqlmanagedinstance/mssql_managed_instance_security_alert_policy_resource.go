// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlManagedInstanceSecurityAlertPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlManagedInstanceSecurityAlertPolicyCreate,
		Read:   resourceMsSqlManagedInstanceSecurityAlertPolicyRead,
		Update: resourceMsSqlManagedInstanceSecurityAlertPolicyUpdate,
		Delete: resourceMsSqlManagedInstanceSecurityAlertPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedInstancesSecurityAlertPolicyID(id)
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

			"managed_instance_name": {
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
						"Brute_Force",
					}, false),
				},
			},

			"email_account_admins_enabled": {
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

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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

func resourceMsSqlManagedInstanceSecurityAlertPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for managed instance security alert policy creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	managedInstanceName := d.Get("managed_instance_name").(string)

	alertPolicy := expandManagedServerSecurityAlertPolicy(d)

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, managedInstanceName, *alertPolicy)
	if err != nil {
		return fmt.Errorf("updataing managed instance security alert policy: %v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creating of managed instance security alert policy (server %q, resource group %q): %+v", managedInstanceName, resourceGroupName, err)
	}

	result, err := client.Get(ctx, resourceGroupName, managedInstanceName)
	if err != nil {
		return fmt.Errorf("retrieving mssql manged instance security alert policy (managed instance %q, resource group %q): %+v", managedInstanceName, resourceGroupName, err)
	}

	if result.Name == nil {
		return fmt.Errorf("reading mssql manged instance security alert policy name (managed instance %q, resource group %q)", managedInstanceName, resourceGroupName)
	}

	id := parse.NewManagedInstancesSecurityAlertPolicyID(subscriptionId, resourceGroupName, managedInstanceName, *result.Name)

	d.SetId(id.ID())

	return resourceMsSqlManagedInstanceSecurityAlertPolicyRead(d, meta)
}

func resourceMsSqlManagedInstanceSecurityAlertPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	instanceName := d.Get("managed_instance_name").(string)

	id, err := parse.ManagedInstancesSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, resourceGroupName, instanceName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	props := existing.SecurityAlertsPolicyProperties
	if props == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("disabled_alerts") {
		if v, ok := d.GetOk("disabled_alerts"); ok {
			disabledAlerts := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				disabledAlerts = append(disabledAlerts, v.(string))
			}
			props.DisabledAlerts = &disabledAlerts
		} else {
			props.DisabledAlerts = nil
		}
	}
	if d.HasChange("email_account_admins_enabled") {
		props.EmailAccountAdmins = utils.Bool(d.Get("email_account_admins_enabled").(bool))
	}

	if d.HasChange("retention_days") {
		props.RetentionDays = utils.Int32(int32(d.Get("retention_days").(int)))
	}

	if d.HasChange("email_addresses") {
		if v, ok := d.GetOk("email_addresses"); ok {
			emailAddresses := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				emailAddresses = append(emailAddresses, v.(string))
			}
			props.EmailAddresses = &emailAddresses
		} else {
			props.EmailAddresses = nil
		}
	}

	if d.HasChange("enabled") {
		if d.Get("enabled").(bool) {
			props.State = sql.SecurityAlertsPolicyStateEnabled
		} else {
			props.State = sql.SecurityAlertsPolicyStateDisabled
		}
	}

	props.StorageAccountAccessKey = utils.String(d.Get("storage_account_access_key").(string))

	if d.HasChange("storage_endpoint") {
		props.StorageEndpoint = utils.String(d.Get("storage_endpoint").(string))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, instanceName, existing)
	if err != nil {
		return fmt.Errorf("updataing managed instance security alert policy: %v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for updating of managed instance security alert policy (server %q, resource group %q): %+v", instanceName, resourceGroupName, err)
	}

	d.SetId(id.ID())

	return resourceMsSqlManagedInstanceSecurityAlertPolicyRead(d, meta)
}

func resourceMsSqlManagedInstanceSecurityAlertPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] reading managed instance security alert policy")

	id, err := parse.ManagedInstancesSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	result, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[WARN] managed instance security alert policy %v not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request to managed instance security alert policy: %+v", err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("managed_instance_name", id.ManagedInstanceName)

	if props := result.SecurityAlertsPolicyProperties; props != nil {
		d.Set("enabled", props.State == sql.SecurityAlertsPolicyStateEnabled)

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
			d.Set("email_account_admins_enabled", props.EmailAccountAdmins)
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

func resourceMsSqlManagedInstanceSecurityAlertPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstancesSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	disabledPolicy := sql.ManagedServerSecurityAlertPolicy{
		SecurityAlertsPolicyProperties: &sql.SecurityAlertsPolicyProperties{
			State: sql.SecurityAlertsPolicyStateDisabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, disabledPolicy)
	if err != nil {
		return fmt.Errorf("updataing managed instance security alert policy: %v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of managed instance security alert policy (server %q, resource group %q): %+v", id.ManagedInstanceName, id.ResourceGroup, err)
	}

	if _, err = client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName); err != nil {
		return fmt.Errorf("deleting managed instance security alert policy: %v", err)
	}

	return nil
}

func expandManagedServerSecurityAlertPolicy(d *pluginsdk.ResourceData) *sql.ManagedServerSecurityAlertPolicy {
	state := sql.SecurityAlertsPolicyStateDisabled

	if d.Get("enabled").(bool) {
		state = sql.SecurityAlertsPolicyStateEnabled
	}

	policy := sql.ManagedServerSecurityAlertPolicy{
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

	if v, ok := d.GetOk("email_account_admins_enabled"); ok {
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
