// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	billingClient := meta.(*clients.Client).AppInsights.BillingClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := componentsapis.ParseComponentID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.ComponentsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	component := existing.Model

	oldWorkspaceId, newWorkspaceId := d.GetChange("workspace_id")
	if oldWorkspaceId.(string) != "" && newWorkspaceId.(string) == "" {
		return fmt.Errorf("`workspace_id` cannot be removed after set. If `workspace_id` is not specified but you encounter a diff, this might indicate a Microsoft initiated automatic migration from classic resources to workspace-based resources. If this is the case, please update `workspace_id` in your config file to the new value")
	}

	if d.HasChange("sampling_percentage") {
		component.Properties.SamplingPercentage = pointer.To(d.Get("sampling_percentage").(float64))
	}

	if d.HasChange("ip_masking_enabled") {
		component.Properties.DisableIPMasking = pointer.To(!d.Get("ip_masking_enabled").(bool))
	}

	if d.HasChange("local_authentication_enabled") {
		component.Properties.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))
	}

	if d.HasChange("internet_ingestion_enabled") {
		component.Properties.PublicNetworkAccessForIngestion = pointer.To(componentsapis.PublicNetworkAccessTypeDisabled)
		if d.Get("internet_ingestion_enabled").(bool) {
			component.Properties.PublicNetworkAccessForIngestion = pointer.To(componentsapis.PublicNetworkAccessTypeEnabled)
		}
	}

	if d.HasChange("internet_query_enabled") {
		component.Properties.PublicNetworkAccessForQuery = pointer.To(componentsapis.PublicNetworkAccessTypeDisabled)
		if d.Get("internet_query_enabled").(bool) {
			component.Properties.PublicNetworkAccessForQuery = pointer.To(componentsapis.PublicNetworkAccessTypeEnabled)
		}
	}

	if d.HasChange("force_customer_storage_for_profiler") {
		component.Properties.ForceCustomerStorageForProfiler = pointer.To(d.Get("force_customer_storage_for_profiler").(bool))
	}

	if d.HasChange("workspace_id") {
		workspaceID, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
		if err != nil {
			return err
		}
		component.Properties.WorkspaceResourceId = pointer.To(workspaceID.ID())
	}

	if d.HasChange("retention_in_days") {
		component.Properties.RetentionInDays = pointer.To(int64(d.Get("retention_in_days").(int)))
	}

	if d.HasChange("tags") {
		component.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err = client.ComponentsCreateOrUpdate(ctx, *id, *component); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	read, err := client.ComponentsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if read.Model.Id == nil {
		return fmt.Errorf("retrieving %s: `id` was nil", id)
	}
	billingId, err := componentfeaturesandpricingapis.ParseComponentID(id.ID())
	if err != nil {
		return err
	}
	billingExisting, err := billingClient.ComponentCurrentBillingFeaturesGet(ctx, *billingId)
	if err != nil {
		return fmt.Errorf("retrieving Billing Features for %s: %+v", id, err)
	}

	if billingExisting.Model == nil {
		return fmt.Errorf("retrieving Billing Features for %s: `model` was nil", id)
	}

	billingProps := billingExisting.Model

	if billingProps.DataVolumeCap == nil {
		billingProps.DataVolumeCap = &componentfeaturesandpricingapis.ApplicationInsightsComponentDataVolumeCap{}
	}

	if d.HasChange("daily_data_cap_in_gb") {
		billingProps.DataVolumeCap.Cap = pointer.To(d.Get("daily_data_cap_in_gb").(float64))
	}

	if d.HasChange("daily_data_cap_notifications_enabled") {
		billingProps.DataVolumeCap.StopSendNotificationWhenHitCap = pointer.To(!d.Get("daily_data_cap_notifications_enabled").(bool))
	}

	if _, err = billingClient.ComponentCurrentBillingFeaturesUpdate(ctx, *billingId, *billingProps); err != nil {
		return fmt.Errorf("updating Billing Features for %s: %+v", id, err)
	}

	return resourceApplicationInsightsRead(d, meta)
}
