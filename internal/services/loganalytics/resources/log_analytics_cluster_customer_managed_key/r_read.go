// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_cluster_customer_managed_key

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsClusterCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keyVaultKeyId := ""
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if kvProps := props.KeyVaultProperties; kvProps != nil {
				keyVaultUri := pointer.From(kvProps.KeyVaultUri)
				keyName := pointer.From(kvProps.KeyName)
				keyVersion := pointer.From(kvProps.KeyVersion)

				if keyVaultUri != "" && keyName != "" {
					keyId, err := keyvault.NewNestedItemID(keyVaultUri, keyvault.NestedItemTypeKey, keyName, keyVersion)
					if err != nil {
						return err
					}
					keyVaultKeyId = keyId.ID()
				}
			}
		}
	}

	if keyVaultKeyId == "" {
		log.Printf("[DEBUG] %s has no Customer Managed Key - removing from state", *id)
		d.SetId("")
		return nil
	}

	d.Set("log_analytics_cluster_id", d.Id())
	d.Set("key_vault_key_id", keyVaultKeyId)

	return nil
}
