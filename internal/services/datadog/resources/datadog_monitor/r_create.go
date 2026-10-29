// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogMonitorCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Datadog.DatadogMonitorResources
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := datadogmonitorresources.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.MonitorsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_datadog_monitor", id.ID())
		}
	}

	monitoringStatus := datadogmonitorresources.MonitoringStatusDisabled
	if d.Get("monitoring_enabled").(bool) {
		monitoringStatus = datadogmonitorresources.MonitoringStatusEnabled
	}

	payload := datadogmonitorresources.DatadogMonitorResource{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: expandMonitorIdentityProperties(d.Get("identity").([]any)),
		Sku: &datadogmonitorresources.ResourceSku{
			Name: d.Get("sku_name").(string),
		},
		Properties: &datadogmonitorresources.MonitorProperties{
			DatadogOrganizationProperties: expandMonitorOrganizationProperties(d.Get("datadog_organization").([]any)),
			UserInfo:                      expandMonitorUserInfo(d.Get("user").([]any)),
			MonitoringStatus:              pointer.To(monitoringStatus),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.MonitorsCreateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatadogMonitorRead(d, meta)
}
