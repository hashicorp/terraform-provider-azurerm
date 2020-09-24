package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMSSQLManagedInstanceKeys() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMSSQLManagedInstanceKeysCreateUpdate,
		Read:   resourceArmMSSQLManagedInstanceKeysRead,
		Update: resourceArmMSSQLManagedInstanceKeysCreateUpdate,
		Delete: resourceArmMSSQLManagedInstanceKeysDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"managed_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"server_key_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"uri": {
				Type:         schema.TypeString,
				Optional:	  true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			
			
			"thumbprint": {
				Type:             schema.TypeString,				
				Computed:         true,
			},

			"creation_date": {
				Type:             schema.TypeString,				
				Computed:         true,
			},

			"name": {
				Type:             schema.TypeString,				
				Computed:         true,
			},

			"type": {
				Type:             schema.TypeString,				
				Computed:         true,
			},

			"kind": {
				Type:             schema.TypeString,				
				Computed:         true,
			},
		},
	}
}

func resourceArmMSSQLManagedInstanceKeysCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	keyClient := meta.(*clients.Client).MSSQL.ManagedInstanceKeysClient
	managedInstanceClient := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managedInstanceId := d.Get("managed_instance_id").(string)
	keyName := d.Get("key_name").(string)

	id, err := azure.ParseAzureResourceID(managedInstanceId)
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	managedInstanceName := id.Path["managedInstances"]

	if _, err := managedInstanceClient.Get(ctx, resGroup, managedInstanceName); err != nil {
		return fmt.Errorf("Error reading managed SQL instance %s: %v", managedInstanceName, err)
	}
	
	if d.IsNewResource() {
		existing, err := keyClient.Get(ctx, resGroup, managedInstanceName, keyName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing key %q ( Managed instance %q, Resource Group %q): %+v", keyName, managedInstanceName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_managed_instance_keys", *existing.ID)
		}
	}

	managedInstanceKey := sql.ManagedInstanceKey{
		ManagedInstanceKeyProperties: &sql.ManagedInstanceKeyProperties{
		},
	}

	if v, exists := d.GetOk("login_username"); exists {
		managedInstanceAdmin.ManagedInstanceAdministratorProperties.Login =  utils.String(v.(string))
	}

	if v, exists := d.GetOk("tenant_id"); exists {
		tid, _ := uuid.FromString(v.(string))
		managedInstanceAdmin.ManagedInstanceAdministratorProperties.TenantID = &tid
	}

	adminFuture, err := adminClient.CreateOrUpdate(ctx, resGroup, name, managedInstanceAdmin)
			if err != nil {
				return fmt.Errorf("Error while creating Managed SQL Instance %q AAD admin details (Resource Group %q): %+v", name, resGroup, err)
			}

			if err = adminFuture.WaitForCompletionRef(ctx, adminClient.Client); err != nil {
				return fmt.Errorf("Error while waiting for creation of Managed SQL Instance %q AAD admin details (Resource Group %q): %+v", name, resGroup, err)
			}


			result, err := adminClient.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making get request for managed SQL instance AAD Admin details %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Error getting ID from managed SQL instance %q AAD Admin details (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*result.ID)

	return resourceArmMSSQLManagedInstanceAdminRead(d, meta)

}

func resourceArmMSSQLManagedInstanceKeysRead(d *schema.ResourceData, meta interface{}) error {
	adminClient := meta.(*clients.Client).MSSQL.ManagedInstanceAdministratorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	adminResp, err := adminClient.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading managed instance %s AAD admin: %v", name, err)
	}

	managedInstanceId, _ := azure.GetSQLResourceParentId(d.Id())
	if err != nil {
		return err
	}
	d.Set("managed_instance_id", managedInstanceId)
	d.Set("name", adminResp.Name)
	d.Set("type", adminResp.Type)

	if props := adminResp.ManagedInstanceAdministratorProperties; props != nil {
		d.Set("admin_type", props.AdministratorType)
		d.Set("login_username", props.Login)
		d.Set("object_id", props.Sid.String())
		d.Set("tenant_id", props.TenantID.String())
	}
	return nil
}

func resourceArmMSSQLManagedInstanceKeysDelete(d *schema.ResourceData, meta interface{}) error {
	adminClient := meta.(*clients.Client).MSSQL.ManagedInstanceAdministratorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	future, err := adminClient.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting managed SQL instance %s admin details: %+v", name, err)
	}

	return future.WaitForCompletionRef(ctx, adminClient.Client)
}
