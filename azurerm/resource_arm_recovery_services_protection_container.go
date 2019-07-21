package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesProtectionContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryServicesProtectionContainerCreate,
		Read:   resourceArmRecoveryServicesProtectionContainerRead,
		Update: nil,
		Delete: resourceArmSiteRecoveryProtectionContainerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmRecoveryServicesProtectionContainerCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("recovery_fabric_name").(string)
	name := d.Get("name").(string)

	client := meta.(*ArmClient).recoveryServices.ProtectionContainerClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing recovery services protection container: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group", *existing.ID)
		}
	}

	parameters := siterecovery.CreateProtectionContainerInput{
		Properties: &siterecovery.CreateProtectionContainerInputProperties{},
	}

	future, err := client.Create(ctx, fabricName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating recovery services protection container1: %+v", err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating recovery services protection container2: %+v", err)
	}

	resp, err := client.Get(ctx, fabricName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery protection container3: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmRecoveryServicesProtectionContainerRead(d, meta)
}

func resourceArmRecoveryServicesProtectionContainerRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	name := id.Path["replicationProtectionContainers"]

	client := meta.(*ArmClient).recoveryServices.ProtectionContainerClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, fabricName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on recovery services protection container %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", resp.Location)
	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	d.Set("recovery_fabric_name", fabricName)
	return nil
}

func resourceArmSiteRecoveryProtectionContainerDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	name := id.Path["replicationProtectionContainers"]

	client := meta.(*ArmClient).recoveryServices.ProtectionContainerClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	future, err := client.Delete(ctx, fabricName, name)
	if err != nil {
		return fmt.Errorf("Error deleting recovery services protection container %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of recovery services protection container %q: %+v", name, err)
	}

	return nil
}
