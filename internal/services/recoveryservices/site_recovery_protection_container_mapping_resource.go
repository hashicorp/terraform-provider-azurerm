package recoveryservices

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSiteRecoveryProtectionContainerMapping() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryContainerMappingCreate,
		Read:   resourceSiteRecoveryContainerMappingRead,
		Update: nil,
		Delete: resourceSiteRecoveryServicesContainerMappingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationProtectionContainerMappingsID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
			},
			"recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_replication_policy_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"recovery_source_protection_container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_target_protection_container_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceSiteRecoveryContainerMappingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for presence of existing site recovery protection container mapping %s (fabric %s, container %s): %+v", name, fabricName, protectionContainerName, err)
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
		return fmt.Errorf("creating site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, fabricName, protectionContainerName, name)
	if err != nil {
		return fmt.Errorf("retrieving site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryContainerMappingRead(d, meta)
}

func resourceSiteRecoveryContainerMappingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationProtectionContainerMappingsID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectionContainerMappingName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : %+v", id.String(), err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("recovery_fabric_name", id.ReplicationFabricName)
	d.Set("recovery_source_protection_container_name", resp.Properties.SourceProtectionContainerFriendlyName)
	d.Set("name", resp.Name)
	d.Set("recovery_replication_policy_id", resp.Properties.PolicyID)
	d.Set("recovery_target_protection_container_id", resp.Properties.TargetProtectionContainerID)
	return nil
}

func resourceSiteRecoveryServicesContainerMappingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationProtectionContainerMappingsID(d.Id())
	if err != nil {
		return err
	}

	instanceType := string(siterecovery.InstanceTypeBasicReplicationProviderSpecificContainerMappingInputInstanceTypeA2A)

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	input := siterecovery.RemoveProtectionContainerMappingInput{
		Properties: &siterecovery.RemoveProtectionContainerMappingInputProperties{
			ProviderSpecificInput: &siterecovery.ReplicationProviderContainerUnmappingInput{
				InstanceType: &instanceType,
			},
		},
	}

	future, err := client.Delete(ctx, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectionContainerMappingName, input)
	if err != nil {
		return fmt.Errorf("deleting site recovery protection container mapping %s : %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of site recovery protection container mapping %s : %+v", id.String(), err)
	}

	return nil
}
