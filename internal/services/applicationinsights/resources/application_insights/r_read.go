// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentfeaturesandpricingapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	billingClient := meta.(*clients.Client).AppInsights.BillingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := componentsapis.ParseComponentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComponentsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	billingId, err := componentfeaturesandpricingapis.ParseComponentID(id.ID())
	if err != nil {
		return err
	}
	billingResp, err := billingClient.ComponentCurrentBillingFeaturesGet(ctx, *billingId)
	if err != nil {
		return fmt.Errorf("retrieving Billing Features for %s: %+v", id, err)
	}

	d.Set("name", id.ComponentName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}

		if props := model.Properties; props != nil {
			vals := map[string]string{
				"web":   "web",
				"other": "other",
			}

			if v, ok := vals[strings.ToLower(string(props.ApplicationType))]; ok {
				d.Set("application_type", v)
			} else {
				d.Set("application_type", string(props.ApplicationType))
			}
			d.Set("app_id", props.AppId)
			d.Set("instrumentation_key", props.InstrumentationKey)
			d.Set("sampling_percentage", props.SamplingPercentage)
			d.Set("ip_masking_enabled", !pointer.From(props.DisableIPMasking))
			d.Set("connection_string", props.ConnectionString)
			d.Set("local_authentication_enabled", !pointer.From(props.DisableLocalAuth))
			d.Set("internet_ingestion_enabled", pointer.From(props.PublicNetworkAccessForIngestion) == componentsapis.PublicNetworkAccessTypeEnabled)
			d.Set("internet_query_enabled", pointer.From(props.PublicNetworkAccessForQuery) == componentsapis.PublicNetworkAccessTypeEnabled)
			d.Set("force_customer_storage_for_profiler", props.ForceCustomerStorageForProfiler)
			d.Set("retention_in_days", pointer.From(props.RetentionInDays))
			workspaceId := ""
			if v := props.WorkspaceResourceId; v != nil {
				id, err := workspaces.ParseWorkspaceIDInsensitively(*v)
				if err != nil {
					return err
				}
				workspaceId = id.ID()
			}
			d.Set("workspace_id", workspaceId)
		}
	}

	if model := billingResp.Model; model != nil {
		if props := model.DataVolumeCap; props != nil {
			d.Set("daily_data_cap_in_gb", props.Cap)
			d.Set("daily_data_cap_notifications_enabled", !pointer.From(props.StopSendNotificationWhenHitCap))
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
