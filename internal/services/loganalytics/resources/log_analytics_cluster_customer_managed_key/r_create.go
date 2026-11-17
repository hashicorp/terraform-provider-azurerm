// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_cluster_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsClusterCustomerManagedKeyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Get("log_analytics_cluster_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

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

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		if props.KeyVaultProperties != nil {
			if keyProps := *props.KeyVaultProperties; keyProps.KeyName != nil && *keyProps.KeyName != "" {
				return tf.ImportAsExistsError("azurerm_log_analytics_cluster_customer_managed_key", id.ID())
			}
		}
	}

	// Ensure `associatedWorkspaces` is not present in request, this is a read only property and cannot be sent to the API
	// Error: updating Customer Managed Key for Cluster
	//		performing CreateOrUpdate: unexpected status 400 (400 Bad Request) with error:
	//		InvalidParameter: 'properties.associatedWorkspaces' is a read only property and cannot be set.
	//		Please refer to https://docs.microsoft.com/azure/azure-monitor/log-query/logs-dedicated-clusters#link-a-workspace-to-the-cluster for more information on how to associate a workspace to the cluster.
	props.AssociatedWorkspaces = nil

	keyId, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
		KeyVaultUri: pointer.To(keyId.KeyVaultBaseURL),
		KeyName:     pointer.To(keyId.Name),
		KeyVersion:  pointer.To(keyId.Version),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *id, *model, sdk.SetIDCallback(meta, id, d)); err != nil {
		return fmt.Errorf("creating Customer Managed Key for %s: %+v", *id, err)
	}
	d.SetId(id.ID())

	updateWait, err := logAnalyticsClusterWaitForState(ctx, client, *id)
	if err != nil {
		return err
	}
	if _, err := updateWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish adding Customer Managed Key: %+v", *id, err)
	}

	return resourceLogAnalyticsClusterCustomerManagedKeyRead(d, meta)
}
