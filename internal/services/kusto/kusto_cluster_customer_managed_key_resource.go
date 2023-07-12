// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
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

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoClusterCustomerManagedKeyV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
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
				ValidateFunc: clusters.ValidateClusterID,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequired(commonids.KeyVaultId{}),

			"key_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"key_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
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
	clusterID, err := clusters.ParseClusterID(clusterIDRaw)
	if err != nil {
		return err
	}

	locks.ByName(clusterID.ClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(clusterID.ClusterName, "azurerm_kusto_cluster")

	cluster, err := clusterClient.Get(ctx, *clusterID)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.ClusterName, clusterID.ResourceGroupName, err)
	}
	if cluster.Model.Properties == nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): `ClusterProperties` was nil", clusterID.ClusterName, clusterID.ResourceGroupName)
	}

	// since we're mutating the kusto cluster here, we can use that as the ID
	resourceID := clusterIDRaw

	if d.IsNewResource() {
		// whilst this looks superflurious given encryption is enabled by default, due to the way
		// the Azure API works this technically can be nil
		if cluster.Model.Properties.KeyVaultProperties != nil {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_customer_managed_key", resourceID)
		}
	}

	keyVaultIDRaw := d.Get("key_vault_id").(string)
	keyVaultID, err := commonids.ParseKeyVaultID(keyVaultIDRaw)
	if err != nil {
		return err
	}

	keyVault, err := vaultsClient.Get(ctx, *keyVaultID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *keyVaultID, err)
	}

	softDeleteEnabled := false
	purgeProtectionEnabled := false
	if model := keyVault.Model; model != nil {
		if esd := model.Properties.EnableSoftDelete; esd != nil {
			softDeleteEnabled = *esd
		}
		if epp := model.Properties.EnablePurgeProtection; epp != nil {
			purgeProtectionEnabled = *epp
		}
	}
	if !softDeleteEnabled || !purgeProtectionEnabled {
		return fmt.Errorf("%s must be configured for both Purge Protection and Soft Delete", *keyVaultID)
	}

	keyVaultBaseURL, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultID)
	if err != nil {
		return fmt.Errorf("looking up Key Vault URI from %s: %+v", *keyVaultID, err)
	}

	keyName := d.Get("key_name").(string)
	keyVersion := d.Get("key_version").(string)
	props := clusters.ClusterUpdate{
		Properties: &clusters.ClusterProperties{
			KeyVaultProperties: &clusters.KeyVaultProperties{
				KeyName:     utils.String(keyName),
				KeyVersion:  utils.String(keyVersion),
				KeyVaultUri: utils.String(*keyVaultBaseURL),
			},
		},
	}

	if v, ok := d.GetOk("user_identity"); ok {
		props.Properties.KeyVaultProperties.UserIdentity = utils.String(v.(string))
	}

	err = clusterClient.UpdateThenPoll(ctx, *clusterID, props, clusters.UpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("updating Customer Managed Key for Kusto Cluster %q (Resource Group %q): %+v", clusterID.ClusterName, clusterID.ResourceGroupName, err)
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

	clusterID, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster, err := clusterClient.Get(ctx, *clusterID)
	if err != nil {
		if !response.WasNotFound(cluster.HttpResponse) {
			log.Printf("[DEBUG] Kusto Cluster %q could not be found in Resource Group %q - removing from state!", clusterID.ClusterName, clusterID.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.ClusterName, clusterID.ResourceGroupName, err)
	}
	if cluster.Model == nil || cluster.Model.Properties == nil {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): `ClusterProperties` was nil", clusterID.ClusterName, clusterID.ResourceGroupName)
	}
	if cluster.Model.Properties.KeyVaultProperties == nil {
		log.Printf("[DEBUG] Customer Managed Key was not defined for Kusto Cluster %q (Resource Group %q) - removing from state!", clusterID.ClusterName, clusterID.ResourceGroupName)
		d.SetId("")
		return nil
	}

	props := cluster.Model.Properties.KeyVaultProperties

	keyName := ""
	keyVaultURI := ""
	keyVersion := ""
	userIdentity := ""
	if props != nil {
		if props.KeyName != nil {
			keyName = *props.KeyName
		}
		if props.KeyVaultUri != nil {
			keyVaultURI = *props.KeyVaultUri
		}
		if props.KeyVersion != nil {
			keyVersion = *props.KeyVersion
		}
		if props.UserIdentity != nil {
			userIdentity = *props.UserIdentity
		}
	}

	if keyVaultURI == "" {
		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): `properties.keyVaultProperties.keyVaultUri` was nil", clusterID.ClusterName, clusterID.ResourceGroupName)
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

	clusterID, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(clusterID.ClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(clusterID.ClusterName, "azurerm_kusto_cluster")

	// confirm it still exists prior to trying to update it, else we'll get an error
	cluster, err := client.Get(ctx, *clusterID)
	if err != nil {
		if response.WasNotFound(cluster.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.ClusterName, clusterID.ResourceGroupName, err)
	}

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Kusto Cluster I am
	// making it reset the Kusto cluster to its default state
	props := clusters.ClusterUpdate{
		Properties: &clusters.ClusterProperties{
			KeyVaultProperties: &clusters.KeyVaultProperties{},
		},
	}

	err = client.UpdateThenPoll(ctx, *clusterID, props, clusters.UpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("removing Customer Managed Key for Kusto Cluster %q (Resource Group %q): %+v", clusterID.ClusterName, clusterID.ResourceGroupName, err)
	}

	return nil
}
