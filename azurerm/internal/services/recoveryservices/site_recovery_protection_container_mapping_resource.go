package recoveryservices

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/recoveryservices/validate"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSiteRecoveryProtectionContainerMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRecoveryContainerMappingCreate,
		Read:   resourceSiteRecoveryContainerMappingRead,
		Update: nil,
		Delete: resourceSiteRecoveryServicesContainerMappingDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: validate.RecoveryServicesVaultName,
			},
			"recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_replication_policy_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"recovery_source_protection_container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_target_protection_container_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceSiteRecoveryContainerMappingCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("recovery_fabric_name").(string)
	policyId := d.Get("recovery_replication_policy_id").(string)
	protectionContainerName := d.Get("recovery_source_protection_container_name").(string)
	targetContainerId := d.Get("recovery_target_protection_container_id").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, protectionContainerName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing site recovery protection container mapping %s (fabric %s, container %s): %+v", name, fabricName, protectionContainerName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_protection_container_mapping", handleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	parameters := siterecovery.CreateProtectionContainerMappingInput{
		Properties: &siterecovery.CreateProtectionContainerMappingInputProperties{
			TargetProtectionContainerID: &targetContainerId,
			PolicyID:                    &policyId,
			ProviderSpecificInput:       siterecovery.ReplicationProviderSpecificContainerMappingInput{},
		},
	}
	future, err := client.Create(ctx, fabricName, protectionContainerName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, fabricName, protectionContainerName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryContainerMappingRead(d, meta)
}

func resourceSiteRecoveryContainerMappingRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectionContainerMappings"]

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, fabricName, protectionContainerName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}

	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	d.Set("recovery_fabric_name", fabricName)
	d.Set("recovery_source_protection_container_name", resp.Properties.SourceProtectionContainerFriendlyName)
	d.Set("name", resp.Name)
	d.Set("recovery_replication_policy_id", resp.Properties.PolicyID)
	d.Set("recovery_target_protection_container_id", resp.Properties.TargetProtectionContainerID)
	return nil
}

func resourceSiteRecoveryServicesContainerMappingDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectionContainerMappings"]
	instanceType := string(siterecovery.InstanceTypeBasicReplicationProviderSpecificContainerMappingInputInstanceTypeA2A)

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	input := siterecovery.RemoveProtectionContainerMappingInput{
		Properties: &siterecovery.RemoveProtectionContainerMappingInputProperties{
			ProviderSpecificInput: &siterecovery.ReplicationProviderContainerUnmappingInput{
				InstanceType: &instanceType,
			},
		},
	}

	future, err := client.Delete(ctx, fabricName, protectionContainerName, name, input)
	if err != nil {
		return fmt.Errorf("Error deleting site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}

	return nil
}
