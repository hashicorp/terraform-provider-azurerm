package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesNetworkMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryNetworkMappingCreate,
		Read:   resourceArmRecoveryNetworkMappingRead,
		Update: nil,
		Delete: resourceArmRecoveryNetworkMappingDelete,
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
			"source_recovery_fabric_name": {
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
			"target_recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"source_network_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"source_network_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_network_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmRecoveryNetworkMappingCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	targetFabricName := d.Get("target_recovery_fabric_name").(string)
	sourceNetworkName := d.Get("source_network_name").(string)
	sourceNetworkId := d.Get("source_network_id").(string)
	targetNetworkId := d.Get("target_network_id").(string)
	name := d.Get("name").(string)

	client := meta.(*ArmClient).recoveryServices.NetworkMappingClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	var parameters = siterecovery.CreateNetworkMappingInput{
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
		return fmt.Errorf("Error creating recovery network mapping: %+v", err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating recovery network mapping: %+v", err)
	}

	resp, err := client.Get(ctx, fabricName, sourceNetworkName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving recovery network mapping: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmRecoveryNetworkMappingRead(d, meta)
}

func resourceArmRecoveryNetworkMappingRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	networkName := id.Path["replicationNetworks"]
	name := id.Path["replicationNetworkMappings"]

	client := meta.(*ArmClient).recoveryServices.NetworkMappingClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, fabricName, networkName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on recovery services protection container mapping %q: %+v", name, err)
	}

	targetFabricId, err := parseAzureResourceID(*resp.Properties.RecoveryFabricArmID)
	if err != nil {
		return err
	}

	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	d.Set("source_recovery_fabric_name", fabricName)
	d.Set("source_network_name", resp.Properties.PrimaryNetworkFriendlyName)
	d.Set("name", resp.Name)
	d.Set("target_recovery_fabric_name", targetFabricId.Path["replicationFabrics"])
	d.Set("source_network_id", unifyIdCasing(resp.Properties.PrimaryNetworkID))
	d.Set("target_network_id", unifyIdCasing(resp.Properties.RecoveryNetworkID))

	return nil
}

// So the casing of the IDs are outputted differently from different services in Azure, try to handle that
func unifyIdCasing(id *string) *string {
	if id == nil {
		return id
	}

	s := *id
	s = strings.Replace(s, "/resourcegroups/", "/resourceGroups/", 1)
	s = strings.Replace(s, "/microsoft.network/virtualnetworks/", "/Microsoft.Network/virtualNetworks/", 1)
	return &s
}

func resourceArmRecoveryNetworkMappingDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	networkName := id.Path["replicationNetworks"]
	name := id.Path["replicationNetworkMappings"]

	client := meta.(*ArmClient).recoveryServices.NetworkMappingClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	future, err := client.Delete(ctx, fabricName, networkName, name)
	if err != nil {
		return fmt.Errorf("Error deleting recovery services protection container mapping %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of recovery services protection container mapping %q: %+v", name, err)
	}

	return nil
}
