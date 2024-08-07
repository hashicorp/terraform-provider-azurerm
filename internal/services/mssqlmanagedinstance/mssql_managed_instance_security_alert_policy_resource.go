// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedserversecurityalertpolicies"
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

func resourceMsSqlManagedInstanceSecurityAlertPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstancesSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	existing, err := client.Get(ctx, managedInstanceId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	payload := existing.Model

	if d.HasChange("disabled_alerts") {
		if v, ok := d.GetOk("disabled_alerts"); ok {
			disabledAlerts := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				disabledAlerts = append(disabledAlerts, v.(string))
			}
			payload.Properties.DisabledAlerts = &disabledAlerts
		} else {
			payload.Properties.DisabledAlerts = nil
		}
	}
	if d.HasChange("email_account_admins_enabled") {
		payload.Properties.EmailAccountAdmins = utils.Bool(d.Get("email_account_admins_enabled").(bool))
	}

	if d.HasChange("retention_days") {
		payload.Properties.RetentionDays = pointer.To(int64(d.Get("retention_days").(int)))
	}

	if d.HasChange("email_addresses") {
		if v, ok := d.GetOk("email_addresses"); ok {
			emailAddresses := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				emailAddresses = append(emailAddresses, v.(string))
			}
			payload.Properties.EmailAddresses = &emailAddresses
		} else {
			payload.Properties.EmailAddresses = nil
		}
	}

	if d.HasChange("enabled") {
		if d.Get("enabled").(bool) {
			payload.Properties.State = managedserversecurityalertpolicies.SecurityAlertsPolicyStateEnabled
		} else {
			payload.Properties.State = managedserversecurityalertpolicies.SecurityAlertsPolicyStateDisabled
		}
	}

	if d.HasChange("storage_account_access_key") {
		payload.Properties.StorageAccountAccessKey = utils.String(d.Get("storage_account_access_key").(string))
	}

	// StorageAccountAccessKey cannot be passed in if it is empty. The api returns this as empty so we need to nil it before sending it back to the api
	if payload.Properties.StorageAccountAccessKey != nil && *payload.Properties.StorageAccountAccessKey == "" {
		payload.Properties.StorageAccountAccessKey = nil
	}

	if d.HasChange("storage_endpoint") {
		payload.Properties.StorageEndpoint = utils.String(d.Get("storage_endpoint").(string))
	}

	// StorageEndpoint cannot be passed in if it is empty. The api returns this as empty so we need to nil it before sending it back to the api
	if payload.Properties.StorageEndpoint != nil && *payload.Properties.StorageEndpoint == "" {
		payload.Properties.StorageEndpoint = nil
	}

	err = client.CreateOrUpdateThenPoll(ctx, managedInstanceId, *payload)
	if err != nil {
		return fmt.Errorf("updating managed instance security alert policy: %v", err)
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

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	result, err := client.Get(ctx, managedInstanceId)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			log.Printf("[WARN] managed instance security alert policy %v not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request to managed instance security alert policy: %+v", err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("managed_instance_name", id.ManagedInstanceName)

	if result.Model != nil {

		if props := result.Model.Properties; props != nil {
			d.Set("enabled", props.State == managedserversecurityalertpolicies.SecurityAlertsPolicyStateEnabled)

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

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	disabledPolicy := managedserversecurityalertpolicies.ManagedServerSecurityAlertPolicy{
		Properties: &managedserversecurityalertpolicies.SecurityAlertsPolicyProperties{
			State: managedserversecurityalertpolicies.SecurityAlertsPolicyStateDisabled,
		},
	}

	err = client.CreateOrUpdateThenPoll(ctx, managedInstanceId, disabledPolicy)
	if err != nil {
		return fmt.Errorf("updating managed instance security alert policy: %v", err)
	}

	if _, err = client.Get(ctx, managedInstanceId); err != nil {
		return fmt.Errorf("deleting managed instance security alert policy: %v", err)
	}

	return nil
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
		props.EmailAccountAdmins = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("retention_days"); ok {
		props.RetentionDays = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		props.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_endpoint"); ok {
		props.StorageEndpoint = utils.String(v.(string))
	}

	return &policy
}
