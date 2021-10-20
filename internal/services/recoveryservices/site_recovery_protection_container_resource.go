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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSiteRecoveryProtectionContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryProtectionContainerCreate,
		Read:   resourceSiteRecoveryProtectionContainerRead,
		Update: nil,
		Delete: resourceSiteRecoveryProtectionContainerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationProtectionContainerID(id)
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
		},
	}
}

func resourceSiteRecoveryProtectionContainerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("recovery_fabric_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ProtectionContainerClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing site recovery protection container %s (fabric %s): %+v", name, fabricName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_protection_container", handleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	parameters := siterecovery.CreateProtectionContainerInput{
		Properties: &siterecovery.CreateProtectionContainerInputProperties{},
	}

	future, err := client.Create(ctx, fabricName, name, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery protection container %s (fabric %s): %+v", name, fabricName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating site recovery protection container %s (fabric %s): %+v", name, fabricName, err)
	}

	resp, err := client.Get(ctx, fabricName, name)
	if err != nil {
		return fmt.Errorf("retrieving site recovery protection container %s (fabric %s): %+v", name, fabricName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryProtectionContainerRead(d, meta)
}

func resourceSiteRecoveryProtectionContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ProtectionContainerClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id.ReplicationFabricName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery protection container %s : %+v", id.String(), err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("recovery_fabric_name", id.ReplicationFabricName)
	return nil
}

func resourceSiteRecoveryProtectionContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ProtectionContainerClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, id.ReplicationFabricName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting site recovery protection container %s : %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of site recovery protection container %s : %+v", id.String(), err)
	}

	return nil
}
