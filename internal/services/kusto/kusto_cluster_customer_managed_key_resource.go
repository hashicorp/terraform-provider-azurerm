package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoClusterCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterCustomerManagedKeyCreateUpdate,
		Read:   resourceKustoClusterCustomerManagedKeyRead,
		Update: resourceKustoClusterCustomerManagedKeyCreateUpdate,
		Delete: resourceKustoClusterCustomerManagedKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterID,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"key_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"key_version": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"user_identity": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},
		},
	}
}

func resourceKustoClusterCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): `ClusterProperties` was nil", clusterID.Name, clusterID.ResourceGroup)
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
		return fmt.Errorf("retrieving Key Vault %q (Resource Group %q): %+v", keyVaultID.Name, keyVaultID.ResourceGroup, err)
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
		return fmt.Errorf("looking up Key Vault URI from Key Vault %q (Resource Group %q): %+v", keyVaultID.Name, keyVaultID.ResourceGroup, err)
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

	if v, ok := d.GetOk("user_identity"); ok {
		props.ClusterProperties.KeyVaultProperties.UserIdentity = utils.String(v.(string))
	}

	future, err := clusterClient.Update(ctx, clusterID.ResourceGroup, clusterID.Name, props, "")
	if err != nil {
		return fmt.Errorf("updating Customer Managed Key for Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
		return fmt.Errorf("waiting for completion of Kusto Cluster Update %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}

	d.SetId(resourceID)

	return resourceKustoClusterCustomerManagedKeyRead(d, meta)
}

func resourceKustoClusterCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clusterClient := meta.(*clients.Client).Kusto.ClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
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

		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): `ClusterProperties` was nil", clusterID.Name, clusterID.ResourceGroup)
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
	userIdentity := ""
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
		if props.UserIdentity != nil {
			userIdentity = *props.UserIdentity
		}
	}

	if keyVaultURI == "" {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): `properties.keyVaultProperties.keyVaultUri` was nil", clusterID.Name, clusterID.ResourceGroup)
	}

	// now we have the key vault uri we can look up the ID
	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultURI)
	if err != nil {
		return fmt.Errorf("retrieving Key Vault ID from the Base URI %q: %+v", keyVaultURI, err)
	}

	d.Set("cluster_id", d.Id())
	d.Set("key_vault_id", keyVaultID)
	d.Set("key_name", keyName)
	d.Set("key_version", keyVersion)
	d.Set("user_identity", userIdentity)
	return nil
}

func resourceKustoClusterCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
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

	future, err := client.Update(ctx, clusterID.ResourceGroup, clusterID.Name, props, "")
	if err != nil {
		return fmt.Errorf("removing Customer Managed Key for Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Kusto Cluster Update %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
	}

	return nil
}
