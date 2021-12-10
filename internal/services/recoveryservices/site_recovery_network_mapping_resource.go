package recoveryservices

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSiteRecoveryNetworkMapping() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryNetworkMappingCreate,
		Read:   resourceSiteRecoveryNetworkMappingRead,
		Delete: resourceSiteRecoveryNetworkMappingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationNetworkMappingID(id)
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
			"source_recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"target_recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"source_network_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_network_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceSiteRecoveryNetworkMappingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
	parsedSourceNetworkId, err := networkParse.VirtualNetworkID(sourceNetworkId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_network_id '%s' (network mapping %s): %+v", sourceNetworkId, name, err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, parsedSourceNetworkId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) &&
				// todo this workaround can be removed when this bug is fixed
				// https://github.com/Azure/azure-sdk-for-go/issues/8705
				!utils.ResponseWasStatusCode(existing.Response, http.StatusBadRequest) {
				return fmt.Errorf("checking for presence of existing site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_network_mapping", handleAzureSdkForGoBug2824(*existing.ID))
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
	future, err := client.Create(ctx, fabricName, parsedSourceNetworkId.Name, name, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, fabricName, parsedSourceNetworkId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryNetworkMappingRead(d, meta)
}

func resourceSiteRecoveryNetworkMappingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationNetworkMappingID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id.ReplicationFabricName, id.ReplicationNetworkName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery network mapping %s (vault %s): %+v", id.Name, id.VaultName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("source_recovery_fabric_name", id.ReplicationFabricName)
	d.Set("name", resp.Name)
	if props := resp.Properties; props != nil {
		d.Set("source_network_id", props.PrimaryNetworkID)
		d.Set("target_network_id", props.RecoveryNetworkID)

		targetFabricId, err := parse.ReplicationFabricID(handleAzureSdkForGoBug2824(*resp.Properties.RecoveryFabricArmID))
		if err != nil {
			return err
		}
		d.Set("target_recovery_fabric_name", targetFabricId.Name)
	}

	return nil
}

func resourceSiteRecoveryNetworkMappingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationNetworkMappingID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, id.ReplicationFabricName, id.ReplicationNetworkName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting site recovery network mapping %s (vault %s): %+v", id.Name, id.VaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of site recovery network mapping  %s (vault %s): %+v", id.Name, id.VaultName, err)
	}

	return nil
}
