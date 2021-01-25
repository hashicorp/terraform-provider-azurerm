package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKustoClusterCustomerManagedKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKustoClusterCustomerManagedKeyCreateUpdate,
		Read:   resourceKustoClusterCustomerManagedKeyRead,
		Update: resourceKustoClusterCustomerManagedKeyCreateUpdate,
		Delete: resourceKustoClusterCustomerManagedKeyDelete,

		// TODO: this needs a custom ID validating importer
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
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"key_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceKustoClusterCustomerManagedKeyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterClient := meta.(*clients.Client).Kusto.ClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	vaultsClient := keyVaultsClient.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterIDRaw := d.Get("cluster_id").(string)
	clusterID, err := parse.ClusterID(clusterIDRaw)
	if err != nil {
		return err
	}

	locks.ByName(clusterID.Name, "azurerm_kusto_cluster")
	defer locks.UnlockByName(clusterID.Name, "azurerm_kusto_cluster")

	cluster, err := clusterClient.Get(ctx, clusterID.ResourceGroup, clusterID.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil {
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): `ClusterProperties` was nil", clusterID.Name, clusterID.ResourceGroup)
	}

	// since we're mutating the kusto cluster here, we can use that as the ID
	resourceID := clusterIDRaw

	if d.IsNewResource() {
		// whilst this looks superflurious given encryption is enabled by default, due to the way
		// the Azure API works this technically can be nil
		if cluster.ClusterProperties.KeyVaultProperties != nil {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_customer_managed_key", resourceID)
		}
	}

	keyVaultIDRaw := d.Get("key_vault_id").(string)
	keyVaultID, err := keyVaultParse.VaultID(keyVaultIDRaw)
	if err != nil {
		return err
	}

	keyVault, err := vaultsClient.Get(ctx, keyVaultID.ResourceGroup, keyVaultID.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", keyVaultID.Name, keyVaultID.ResourceGroup, err)
	}

	softDeleteEnabled := false
	purgeProtectionEnabled := false
	if props := keyVault.Properties; props != nil {
		if esd := props.EnableSoftDelete; esd != nil {
			softDeleteEnabled = *esd
		}
		if epp := props.EnablePurgeProtection; epp != nil {
			purgeProtectionEnabled = *epp
		}
	}
	if !softDeleteEnabled || !purgeProtectionEnabled {
		return fmt.Errorf("Key Vault %q (Resource Group %q) must be configured for both Purge Protection and Soft Delete", keyVaultID.Name, keyVaultID.ResourceGroup)
	}

	keyVaultBaseURL, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultID)
	if err != nil {
		return fmt.Errorf("Error looking up Key Vault URI from Key Vault %q (Resource Group %q): %+v", keyVaultID.Name, keyVaultID.ResourceGroup, err)
	}

	keyName := d.Get("key_name").(string)
	keyVersion := d.Get("key_version").(string)
	props := kusto.ClusterUpdate{
		ClusterProperties: &kusto.ClusterProperties{
			KeyVaultProperties: &kusto.KeyVaultProperties{
				KeyName:     utils.String(keyName),
				KeyVersion:  utils.String(keyVersion),
				KeyVaultURI: utils.String(*keyVaultBaseURL),
			},
		},
	}

	future, err := clusterClient.Update(ctx, clusterID.ResourceGroup, clusterID.Name, props)
	if err != nil {
		return fmt.Errorf("Error updating Customer Managed Key for Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto Cluster Update %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}

	d.SetId(resourceID)

	return resourceKustoClusterCustomerManagedKeyRead(d, meta)
}

func resourceKustoClusterCustomerManagedKeyRead(d *schema.ResourceData, meta interface{}) error {
	clusterClient := meta.(*clients.Client).Kusto.ClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterID, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster, err := clusterClient.Get(ctx, clusterID.ResourceGroup, clusterID.Name)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			log.Printf("[DEBUG] Kusto Cluster %q could not be found in Resource Group %q - removing from state!", clusterID.Name, clusterID.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil {
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): `ClusterProperties` was nil", clusterID.Name, clusterID.ResourceGroup)
	}
	if cluster.ClusterProperties.KeyVaultProperties == nil {
		log.Printf("[DEBUG] Customer Managed Key was not defined for Kusto Cluster %q (Resource Group %q) - removing from state!", clusterID.Name, clusterID.ResourceGroup)
		d.SetId("")
		return nil
	}

	props := cluster.ClusterProperties.KeyVaultProperties

	keyName := ""
	keyVaultURI := ""
	keyVersion := ""
	if props != nil {
		if props.KeyName != nil {
			keyName = *props.KeyName
		}
		if props.KeyVaultURI != nil {
			keyVaultURI = *props.KeyVaultURI
		}
		if props.KeyVersion != nil {
			keyVersion = *props.KeyVersion
		}
	}

	if keyVaultURI == "" {
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): `properties.keyVaultProperties.keyVaultUri` was nil", clusterID.Name, clusterID.ResourceGroup)
	}

	// now we have the key vault uri we can look up the ID
	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, keyVaultURI)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault ID from the Base URI %q: %+v", keyVaultURI, err)
	}

	d.Set("cluster_id", d.Id())
	d.Set("key_vault_id", keyVaultID)
	d.Set("key_name", keyName)
	d.Set("key_version", keyVersion)

	return nil
}

func resourceKustoClusterCustomerManagedKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterID, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(clusterID.Name, "azurerm_kusto_cluster")
	defer locks.UnlockByName(clusterID.Name, "azurerm_kusto_cluster")

	// confirm it still exists prior to trying to update it, else we'll get an error
	cluster, err := client.Get(ctx, clusterID.ResourceGroup, clusterID.Name)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			return nil
		}

		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Kusto Cluster I am
	// making it reset the Kusto cluster to its default state
	props := kusto.ClusterUpdate{
		ClusterProperties: &kusto.ClusterProperties{
			KeyVaultProperties: &kusto.KeyVaultProperties{},
		},
	}

	future, err := client.Update(ctx, clusterID.ResourceGroup, clusterID.Name, props)
	if err != nil {
		return fmt.Errorf("Error removing Customer Managed Key for Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto Cluster Update %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}

	return nil
}
