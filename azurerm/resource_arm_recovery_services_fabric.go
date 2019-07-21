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

func resourceArmRecoveryServicesFabric() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryServicesFabricCreate,
		Read:   resourceArmRecoveryServicesFabricRead,
		Update: nil,
		Delete: resourceArmRecoveryServicesFabricDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"location": azure.SchemaLocation(),
		},
	}
}

func resourceArmRecoveryServicesFabricCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	name := d.Get("name").(string)

	client := meta.(*ArmClient).recoveryServices.FabricClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing recovery services fabric: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group", *existing.ID)
		}
	}

	parameters := siterecovery.FabricCreationInput{
		Properties: &siterecovery.FabricCreationInputProperties{
			CustomDetails: siterecovery.AzureFabricCreationInput{
				InstanceType: "Azure",
				Location:     &location,
			},
		},
	}

	future, err := client.Create(ctx, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating recovery services fabric: %+v", err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating recovery services fabric: %+v", err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving recovery services fabric: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmRecoveryServicesFabricRead(d, meta)
}

func resourceArmRecoveryServicesFabricRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationFabrics"]

	client := meta.(*ArmClient).recoveryServices.FabricClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on recovery services fabric %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", resp.Properties.FriendlyName) // Crazy? yes. But the location comes back in the firendly name
	d.Set("recovery_vault_name", vaultName)
	return nil
}

func resourceArmRecoveryServicesFabricDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationFabrics"]

	client := meta.(*ArmClient).recoveryServices.FabricClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	future, err := client.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting recovery services fabric %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of recovery services fabric %q: %+v", name, err)
	}

	return nil
}
