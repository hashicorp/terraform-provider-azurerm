// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	clusterClient := meta.(*clients.Client).Kusto.ClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	vaultsClient := keyVaultsClient.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterIDRaw := d.Get("cluster_id").(string)
	clusterID, err := commonids.ParseKustoClusterID(clusterIDRaw)
	if err != nil {
		return err
	}

	locks.ByName(clusterID.KustoClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(clusterID.KustoClusterName, "azurerm_kusto_cluster")

	cluster, err := clusterClient.Get(ctx, *clusterID)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *clusterID, err)
	}
	if cluster.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `ClusterProperties` was nil", *clusterID)
	}

	// since we're mutating the kusto cluster here, we can use that as the ID
	resourceID := clusterIDRaw

	if d.IsNewResource() {
		// whilst this looks superflurious given encryption is enabled by default, due to the way
		// the Azure API works this technically can be nil
		{
			if cluster.Model.Properties.KeyVaultProperties != nil {
				return tf.ImportAsExistsError("azurerm_kusto_cluster_customer_managed_key", resourceID)
			}
		}
	}

	keyName := ""
	keyVersion := ""
	keyVaultURI := ""

	if _, ok := d.GetOk("key_vault_id"); ok {
		keyVaultID, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
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
		keyName = d.Get("key_name").(string)
		keyVersion = d.Get("key_version").(string)
		keyVaultURI = *keyVaultBaseURL
	} else if managedHSMKeyId, ok := d.GetOk("managed_hsm_key_id"); ok {
		if keyId, err := parse.ManagedHSMDataPlaneVersionedKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = keyId.KeyName
			keyVersion = keyId.KeyVersion
			keyVaultURI = keyId.BaseUri()
		} else if keyId, err := parse.ManagedHSMDataPlaneVersionlessKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = keyId.KeyName
			keyVersion = ""
			keyVaultURI = keyId.BaseUri()
		} else {
			return fmt.Errorf("failed to parse '%s' as HSM key ID", managedHSMKeyId)
		}
	}

	props := clusters.ClusterUpdate{
		Properties: &clusters.ClusterProperties{
			KeyVaultProperties: &clusters.KeyVaultProperties{
				KeyName:     pointer.To(keyName),
				KeyVersion:  pointer.To(keyVersion),
				KeyVaultUri: pointer.To(keyVaultURI),
			},
		},
	}

	if v, ok := d.GetOk("user_identity"); ok {
		props.Properties.KeyVaultProperties.UserIdentity = pointer.To(v.(string))
	}

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err = clusterClient.UpdateThenPoll(ctx, *clusterID, props, clusters.UpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating Customer Managed Key for %s: %+v", *clusterID, err)
	}

	d.SetId(resourceID)

	return resourceKustoClusterCustomerManagedKeyRead(d, meta)
}
