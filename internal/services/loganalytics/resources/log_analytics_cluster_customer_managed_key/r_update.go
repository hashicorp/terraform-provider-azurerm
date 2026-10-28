// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_cluster_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsClusterCustomerManagedKeyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	keyId, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return fmt.Errorf("parsing Key Vault Key ID: %+v", err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `model` is nil", *id)
	}

	if props := model.Properties; props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	// This is a read only property, please see comment in the create function.
	model.Properties.AssociatedWorkspaces = nil

	model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
		KeyVaultUri: pointer.To(keyId.KeyVaultBaseURL),
		KeyName:     pointer.To(keyId.Name),
		KeyVersion:  pointer.To(keyId.Version),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
	}

	return resourceLogAnalyticsClusterCustomerManagedKeyRead(d, meta)
}
