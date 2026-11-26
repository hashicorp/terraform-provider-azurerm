// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_cluster_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsClusterCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
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

	if props.KeyVaultProperties == nil {
		return fmt.Errorf("deleting `azurerm_log_analytics_cluster_customer_managed_key` %s: `customer managed key does not exist!`", *id)
	}

	if props.KeyVaultProperties.KeyName == nil || *props.KeyVaultProperties.KeyName == "" {
		return fmt.Errorf("deleting `azurerm_log_analytics_cluster_customer_managed_key` %s: `customer managed key does not exist!`", *id)
	}

	// This is a read only property, please see comment in the create function.
	props.AssociatedWorkspaces = nil

	// The API only removes the CMK when it is sent empty string values, sending nil for each property or an empty object does not work.
	model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
		KeyVaultUri: pointer.To(""),
		KeyName:     pointer.To(""),
		KeyVersion:  pointer.To(""),
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("deleting Customer Managed Key from %s: %+v", *id, err)
	}

	return nil
}
