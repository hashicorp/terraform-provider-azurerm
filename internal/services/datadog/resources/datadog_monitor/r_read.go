// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogMonitorRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.DatadogMonitorResources
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datadogmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Datadog monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if err := d.Set("identity", flattenMonitorIdentityProperties(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		if props := model.Properties; props != nil {
			if err := d.Set("datadog_organization", flattenMonitorOrganizationProperties(props.DatadogOrganizationProperties, d)); err != nil {
				return fmt.Errorf("setting `datadog_organization`: %+v", err)
			}

			monitoringEnabled := false
			if props.MonitoringStatus != nil && *props.MonitoringStatus == datadogmonitorresources.MonitoringStatusEnabled {
				monitoringEnabled = true
			}

			d.Set("monitoring_enabled", monitoringEnabled)
			d.Set("marketplace_subscription_status", pointer.FromEnum(props.MarketplaceSubscriptionStatus))
		}

		skuName := ""
		if model.Sku != nil {
			skuName = model.Sku.Name
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func flattenMonitorIdentityProperties(input *datadogmonitorresources.IdentityProperties) []any {
	if input == nil || input.Type == nil {
		return make([]any, 0)
	}

	var t string
	if *input.Type != "" {
		t = string(*input.Type)
	}
	return []any{
		map[string]any{
			"type":         t,
			"principal_id": pointer.From(input.PrincipalId),
			"tenant_id":    pointer.From(input.TenantId),
		},
	}
}

func flattenMonitorOrganizationProperties(input *datadogmonitorresources.DatadogOrganizationProperties, d *pluginsdk.ResourceData) []any {
	organisationProperties := d.Get("datadog_organization").([]any)
	if len(organisationProperties) == 0 {
		return make([]any, 0)
	}
	v := organisationProperties[0].(map[string]any)

	return []any{
		map[string]any{
			"name":              pointer.From(input.Name),
			"api_key":           pointer.To(v["api_key"].(string)),
			"application_key":   pointer.To(v["application_key"].(string)),
			"enterprise_app_id": pointer.From(input.EnterpriseAppId),
			"linking_auth_code": pointer.To(v["linking_auth_code"].(string)),
			"linking_client_id": pointer.To(v["linking_client_id"].(string)),
			"redirect_uri":      pointer.From(input.RedirectUri),
			"id":                pointer.From(input.Id),
		},
	}
}

func SkuNameDiffSuppress(_, old, new string, _ *pluginsdk.ResourceData) bool {
	// During creating, accepts any sku name.
	if old == "" {
		return false
	}
	// Sku name of the datadog monitor has two kinds:
	// - Concrete sku: E.g. "payg_v2_Monthly". These will be returned unchanged by API.
	// - Linked sku:  The value is named "Linked". This will be accepted and changed by the API, which will then return you the linked concrete sku.
	if new == "Linked" {
		return true
	}
	return old == new
}
