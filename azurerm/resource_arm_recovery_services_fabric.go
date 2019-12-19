package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesFabric() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "`azurerm_recovery_services_fabric` resource is deprecated in favor of `azurerm_site_recovery_fabric` and will be removed in v2.0 of the AzureRM Provider",
		Create:             resourceArmRecoveryServicesFabricCreate,
		Read:               resourceArmRecoveryServicesFabricRead,
		Update:             nil,
		Delete:             resourceArmRecoveryServicesFabricDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
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

	client := meta.(*clients.Client).RecoveryServices.FabricClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing rec	overy services fabric %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group", azure.HandleAzureSdkForGoBug2824(*existing.ID))
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
		return fmt.Errorf("Error creating recovery services fabric %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating recovery services fabric %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving recovery services fabric %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(azure.HandleAzureSdkForGoBug2824(*resp.ID))

	return resourceArmRecoveryServicesFabricRead(d, meta)
}

func resourceArmRecoveryServicesFabricRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationFabrics"]

	client := meta.(*clients.Client).RecoveryServices.FabricClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on recovery services fabric %s (vault %s): %+v", name, vaultName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if props := resp.Properties; props != nil {
		if azureDetails, isAzureDetails := props.CustomDetails.AsAzureFabricSpecificDetails(); isAzureDetails {
			d.Set("location", azureDetails.Location)
		}
	}
	d.Set("recovery_vault_name", vaultName)
	return nil
}

func resourceArmRecoveryServicesFabricDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationFabrics"]

	client := meta.(*clients.Client).RecoveryServices.FabricClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting recovery services fabric %s (vault %s): %+v", name, vaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of recovery services fabric %s (vault %s): %+v", name, vaultName, err)
	}

	return nil
}
