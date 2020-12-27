package recoveryservices

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSiteRecoveryNetworkMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRecoveryNetworkMappingCreate,
		Read:   resourceSiteRecoveryNetworkMappingRead,
		Delete: resourceSiteRecoveryNetworkMappingDelete,
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
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
			},
			"source_recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"target_recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"source_network_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_network_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceSiteRecoveryNetworkMappingCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	targetFabricName := d.Get("target_recovery_fabric_name").(string)
	sourceNetworkId := d.Get("source_network_id").(string)
	targetNetworkId := d.Get("target_network_id").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// get network name from id
	parsedSourceNetworkId, err := azure.ParseAzureResourceID(sourceNetworkId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_network_id '%s' (network mapping %s): %+v", sourceNetworkId, name, err)
	}
	sourceNetworkName, hasName := parsedSourceNetworkId.Path["virtualNetworks"]
	if !hasName {
		sourceNetworkName, hasName = parsedSourceNetworkId.Path["virtualnetworks"] // Handle that different APIs return different ID casings
		if !hasName {
			return fmt.Errorf("[ERROR] parsed source_network_id '%s' doesn't contain 'virtualnetworks'", parsedSourceNetworkId)
		}
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, sourceNetworkName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) &&
				// todo this workaround can be removed when this bug is fixed
				// https://github.com/Azure/azure-sdk-for-go/issues/8705
				!utils.ResponseWasStatusCode(existing.Response, http.StatusBadRequest) {
				return fmt.Errorf("Error checking for presence of existing site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_network_mapping", azure.HandleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	parameters := siterecovery.CreateNetworkMappingInput{
		Properties: &siterecovery.CreateNetworkMappingInputProperties{
			RecoveryNetworkID:  &targetNetworkId,
			RecoveryFabricName: &targetFabricName,
			FabricSpecificDetails: siterecovery.AzureToAzureCreateNetworkMappingInput{
				PrimaryNetworkID: &sourceNetworkId,
			},
		},
	}
	future, err := client.Create(ctx, fabricName, sourceNetworkName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, fabricName, sourceNetworkName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(azure.HandleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryNetworkMappingRead(d, meta)
}

func resourceSiteRecoveryNetworkMappingRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	networkName := id.Path["replicationNetworks"]
	name := id.Path["replicationNetworkMappings"]

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, fabricName, networkName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}

	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	d.Set("source_recovery_fabric_name", fabricName)
	d.Set("name", resp.Name)
	if props := resp.Properties; props != nil {
		d.Set("source_network_id", props.PrimaryNetworkID)
		d.Set("target_network_id", props.RecoveryNetworkID)

		targetFabricId, err := azure.ParseAzureResourceID(azure.HandleAzureSdkForGoBug2824(*resp.Properties.RecoveryFabricArmID))
		if err != nil {
			return err
		}
		d.Set("target_recovery_fabric_name", targetFabricId.Path["replicationFabrics"])
	}

	return nil
}

func resourceSiteRecoveryNetworkMappingDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	networkName := id.Path["replicationNetworks"]
	name := id.Path["replicationNetworkMappings"]

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, fabricName, networkName, name)
	if err != nil {
		return fmt.Errorf("Error deleting site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of site recovery network mapping  %s (vault %s): %+v", name, vaultName, err)
	}

	return nil
}
