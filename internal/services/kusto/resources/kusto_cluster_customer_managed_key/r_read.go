// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_customer_managed_key

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta any) error {
	clusterClient := meta.(*clients.Client).Kusto.ClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	env := meta.(*clients.Client).Account.Environment
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster, err := clusterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(cluster.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if cluster.Model == nil || cluster.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `ClusterProperties` was nil", id)
	}
	if cluster.Model.Properties.KeyVaultProperties == nil {
		log.Printf("[DEBUG] Customer Managed Key was not defined for %s - removing from state!", id)
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
		return fmt.Errorf("retrieving %s: `properties.keyVaultProperties.keyVaultUri` was nil", id)
	}

	isHSMURI, instanceName, domainSuffix, err := helpers.IsManagedHSMURI(env, keyVaultURI)
	if err != nil {
		return err
	}
	switch {
	case isHSMURI && keyVersion == "":
		{
			d.Set("managed_hsm_key_id", parse.NewManagedHSMDataPlaneVersionlessKeyID(instanceName, domainSuffix, keyName).ID())
		}
	case isHSMURI && keyVersion != "":
		{
			d.Set("managed_hsm_key_id", parse.NewManagedHSMDataPlaneVersionedKeyID(instanceName, domainSuffix, keyName, keyVersion).ID())
		}
	case !isHSMURI:
		{
			// now we have the key vault uri we can look up the ID
			subscriptionResourceId := commonids.NewSubscriptionID(id.SubscriptionId)
			keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultURI)
			if err != nil {
				return fmt.Errorf("retrieving Key Vault ID from the Base URI %q: %+v", keyVaultURI, err)
			}
			d.Set("key_vault_id", keyVaultID)
			d.Set("key_name", keyName)
		}
	}

	d.Set("cluster_id", d.Id())
	d.Set("key_version", keyVersion)
	d.Set("user_identity", userIdentity)
	return nil
}
