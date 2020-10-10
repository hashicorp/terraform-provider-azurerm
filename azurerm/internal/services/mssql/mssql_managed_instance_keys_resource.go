package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: azure.ValidateResourceID,
			},

			"uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"server_key_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
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
			return tf.ImportAsExistsError("azurerm_mssql_managed_instance_key", *existing.ID)
		}
	}

	managedInstanceKey := sql.ManagedInstanceKey{
		ManagedInstanceKeyProperties: &sql.ManagedInstanceKeyProperties{
			ServerKeyType: sql.AzureKeyVault,
		},
	}

	if v, exists := d.GetOk("uri"); exists {
		managedInstanceKey.ManagedInstanceKeyProperties.URI = utils.String(v.(string))
	}

	keyFuture, err := keyClient.CreateOrUpdate(ctx, resGroup, managedInstanceName, keyName, managedInstanceKey)
	if err != nil {
		return fmt.Errorf("Error while creating Managed SQL Instance encryption key %q  (Managed instance %q, Resource Group %q): %+v", keyName, managedInstanceName, resGroup, err)
	}

	if err = keyFuture.WaitForCompletionRef(ctx, keyClient.Client); err != nil {
		return fmt.Errorf("Error while waiting for creation of Managed SQL Instance encryption key %q  (Managed instance %q, Resource Group %q): %+v", keyName, managedInstanceName, resGroup, err)
	}

	result, err := keyClient.Get(ctx, resGroup, managedInstanceName, keyName)
	if err != nil {
		return fmt.Errorf("Error making get request for managed SQL instance encryption key %q  (Managed instance %q, Resource Group %q): %+v", keyName, managedInstanceName, resGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Error getting ID from managed SQL instance encryption key %q  (Managed instance %q, Resource Group %q)", keyName, managedInstanceName, resGroup)
	}

	d.SetId(*result.ID)

	return resourceArmMSSQLManagedInstanceKeysRead(d, meta)

}

func resourceArmMSSQLManagedInstanceKeysRead(d *schema.ResourceData, meta interface{}) error {
	keyClient := meta.(*clients.Client).MSSQL.ManagedInstanceKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	managedInstanceName := id.Path["managedInstances"]
	keyName := id.Path["keys"]

	result, err := keyClient.Get(ctx, resGroup, managedInstanceName, keyName)
	if err != nil {
		return fmt.Errorf("Error while fetching managed SQL instance encryption key %q details (Managed instance %q, Resource Group %q): %+v", keyName, managedInstanceName, resGroup, err)
	}

	managedInstanceId, _ := azure.GetSQLResourceParentId(d.Id())
	if err != nil {
		return err
	}
	d.Set("managed_instance_id", managedInstanceId)
	d.Set("key_name", result.Name)
	d.Set("kind", result.Kind)
	d.Set("type", result.Type)
	d.Set("name", result.Name)

	if props := result.ManagedInstanceKeyProperties; props != nil {
		d.Set("server_key_type", string(props.ServerKeyType))
		d.Set("uri", props.URI)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("creation_date", props.CreationDate.String())
	}
	return nil
}

func resourceArmMSSQLManagedInstanceKeysDelete(d *schema.ResourceData, meta interface{}) error {
	keyClient := meta.(*clients.Client).MSSQL.ManagedInstanceKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	managedInstanceName := id.Path["managedInstances"]
	keyName := id.Path["keys"]

	future, err := keyClient.Delete(ctx, resGroup, managedInstanceName, keyName)
	if err != nil {
		return fmt.Errorf("Error deleting managed SQL instance encryption key %q details (Managed instance %q, Resource Group %q): %+v", keyName, managedInstanceName, resGroup, err)
	}

	return future.WaitForCompletionRef(ctx, keyClient.Client)
}
