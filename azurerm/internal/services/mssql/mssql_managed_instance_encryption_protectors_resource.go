package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMSSQLManagedInstanceEncryptionProtector() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMSSQLManagedInstanceEncryptionProtectorCreateUpdate,
		Read:   resourceArmMSSQLManagedInstanceEncryptionProtectorRead,
		Update: resourceArmMSSQLManagedInstanceEncryptionProtectorCreateUpdate,
		Delete: resourceArmMSSQLManagedInstanceEncryptionProtectorResetToDefault,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"managed_instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_key_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotEmpty,
			},

			"server_key_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ServiceManaged),
					string(sql.AzureKeyVault),
				}, false),
			},

			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
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

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmMSSQLManagedInstanceEncryptionProtectorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	encryptionClient := meta.(*clients.Client).MSSQL.ManagedInstanceEncryptionProtectorsClient
	managedInstanceClient := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managedInstanceName := d.Get("managed_instance_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if _, err := managedInstanceClient.Get(ctx, resGroup, managedInstanceName); err != nil {
		return fmt.Errorf("Error reading managed SQL instance %s: %v", managedInstanceName, err)
	}

	keyName := d.Get("server_key_name").(string)
	serverKeyType := d.Get("server_key_type").(string)

	managedInstanceEncryption := sql.ManagedInstanceEncryptionProtector{
		ManagedInstanceEncryptionProtectorProperties: &sql.ManagedInstanceEncryptionProtectorProperties{
			ServerKeyType: sql.ServerKeyType(serverKeyType),
			ServerKeyName: utils.String(keyName),
		},
	}

	encryptionFuture, err := encryptionClient.CreateOrUpdate(ctx, resGroup, managedInstanceName, managedInstanceEncryption)
	if err != nil {
		return fmt.Errorf("Error while creating Managed SQL Instance %q encryption details (Resource Group %q): %+v", managedInstanceName, resGroup, err)
	}

	if err = encryptionFuture.WaitForCompletionRef(ctx, encryptionClient.Client); err != nil {
		return fmt.Errorf("Error while waiting for creation of Managed SQL Instance %q encryption details (Resource Group %q): %+v", managedInstanceName, resGroup, err)
	}

	result, err := encryptionClient.Get(ctx, resGroup, managedInstanceName)
	if err != nil {
		return fmt.Errorf("Error making get request for managed SQL instance encryption details %q (Resource Group %q): %+v", managedInstanceName, resGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Error getting ID from managed SQL instance %q encryption details (Resource Group %q)", managedInstanceName, resGroup)
	}

	d.SetId(*result.ID)

	return resourceArmMSSQLManagedInstanceEncryptionProtectorRead(d, meta)
}

func resourceArmMSSQLManagedInstanceEncryptionProtectorRead(d *schema.ResourceData, meta interface{}) error {
	encryptionClient := meta.(*clients.Client).MSSQL.ManagedInstanceEncryptionProtectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	managedInstanceName := id.Path["managedInstances"]

	resp, err := encryptionClient.Get(ctx, resGroup, managedInstanceName)
	if err != nil {
		return fmt.Errorf("Error reading managed instance %s encryption details (Resource Group %q): %+v", managedInstanceName, resGroup, err)
	}

	d.Set("managed_instance_name", managedInstanceName)
	d.Set("resource_group_name", resGroup)
	d.Set("name", resp.Name)
	d.Set("type", resp.Type)
	d.Set("kind", resp.Kind)

	if props := resp.ManagedInstanceEncryptionProtectorProperties; props != nil {
		d.Set("server_key_name", props.ServerKeyName)
		d.Set("server_key_type", string(props.ServerKeyType))
		d.Set("uri", props.URI)
		d.Set("thumbprint", props.Thumbprint)
	}
	return nil
}

// Managed Instance Does not support encryption protector deletion.
// Therefore the destroy can default back to ServiceManaged Key encryption rather than to any BYOK TDE protector
func resourceArmMSSQLManagedInstanceEncryptionProtectorResetToDefault(d *schema.ResourceData, meta interface{}) error {
	encryptionClient := meta.(*clients.Client).MSSQL.ManagedInstanceEncryptionProtectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	managedInstanceName := id.Path["managedInstances"]

	managedInstanceEncryption := sql.ManagedInstanceEncryptionProtector{
		ManagedInstanceEncryptionProtectorProperties: &sql.ManagedInstanceEncryptionProtectorProperties{
			ServerKeyType: sql.ServiceManaged,
			ServerKeyName: utils.String("ServiceManaged"),
		},
	}

	encryptionFuture, err := encryptionClient.CreateOrUpdate(ctx, resGroup, managedInstanceName, managedInstanceEncryption)
	if err != nil {
		return fmt.Errorf("Error while creating Managed SQL Instance %q encryption details (Resource Group %q): %+v", managedInstanceName, resGroup, err)
	}

	return encryptionFuture.WaitForCompletionRef(ctx, encryptionClient.Client)
}
