package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// todo 3.0 - this may want to be put into the mssql_server resource now that it exists.

func resourceMsSqlServerSecurityAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMsSqlServerSecurityAlertPolicyCreateUpdate,
		Read:   resourceMsSqlServerSecurityAlertPolicyRead,
		Update: resourceMsSqlServerSecurityAlertPolicyCreateUpdate,
		Delete: resourceMsSqlServerSecurityAlertPolicyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServerSecurityAlertPolicyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"disabled_alerts": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"email_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"state": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.SecurityAlertPolicyStateDisabled),
					string(sql.SecurityAlertPolicyStateEnabled),
					string(sql.SecurityAlertPolicyStateNew),
				}, true),
			},

			"storage_account_access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceMsSqlServerSecurityAlertPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for mssql server security alert policy creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	alertPolicy := expandSecurityAlertPolicy(d)

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, serverName, *alertPolicy)
	if err != nil {
		return fmt.Errorf("error updataing mssql server security alert policy: %v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for creation/update of mssql server security alert policy (server %q, resource group %q): %+v", serverName, resourceGroupName, err)
	}

	result, err := client.Get(ctx, resourceGroupName, serverName)
	if err != nil {
		return fmt.Errorf("error retrieving mssql server security alert policy (server %q, resource group %q): %+v", serverName, resourceGroupName, err)
	}

	if result.ID == nil {
		return fmt.Errorf("error reading mssql server security alert policy id (server %q, resource group %q)", serverName, resourceGroupName)
	}

	d.SetId(*result.ID)

	return resourceMsSqlServerSecurityAlertPolicyRead(d, meta)
}

func resourceMsSqlServerSecurityAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("error making read request to mssql server security alert policy: %+v", err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	if props := result.SecurityAlertPolicyProperties; props != nil {
		d.Set("state", string(props.State))

		if props.DisabledAlerts != nil {
			disabledAlerts := schema.NewSet(schema.HashString, []interface{}{})
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
			emailAddresses := schema.NewSet(schema.HashString, []interface{}{})
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

func resourceMsSqlServerSecurityAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] deleting mssql server security alert policy.")

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	disabledPolicy := sql.ServerSecurityAlertPolicy{
		SecurityAlertPolicyProperties: &sql.SecurityAlertPolicyProperties{
			State: sql.SecurityAlertPolicyStateDisabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, disabledPolicy)
	if err != nil {
		return fmt.Errorf("error updataing mssql server security alert policy: %v", err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for creation/update of mssql server security alert policy (server %q, resource group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if _, err = client.Get(ctx, id.ResourceGroup, id.ServerName); err != nil {
		return fmt.Errorf("error deleting mssql server security alert policy: %v", err)
	}

	return nil
}

func expandSecurityAlertPolicy(d *schema.ResourceData) *sql.ServerSecurityAlertPolicy {
	state := sql.SecurityAlertPolicyState(d.Get("state").(string))

	policy := sql.ServerSecurityAlertPolicy{
		SecurityAlertPolicyProperties: &sql.SecurityAlertPolicyProperties{
			State: state,
		},
	}

	props := policy.SecurityAlertPolicyProperties

	if v, ok := d.GetOk("disabled_alerts"); ok {
		disabledAlerts := make([]string, 0)
		for _, v := range v.(*schema.Set).List() {
			disabledAlerts = append(disabledAlerts, v.(string))
		}
		props.DisabledAlerts = &disabledAlerts
	}

	if v, ok := d.GetOk("email_addresses"); ok {
		emailAddresses := make([]string, 0)
		for _, v := range v.(*schema.Set).List() {
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
