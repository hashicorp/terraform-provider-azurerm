// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_dashboard

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePortalDashboardCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dashboard.NewDashboardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_portal_dashboard", id.ID())
			}
		}
	}

	props := dashboard.Dashboard{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	var dashboardProperties dashboard.DashboardProperties

	dashboardPropsRaw := d.Get("dashboard_properties").(string)
	if err := json.Unmarshal([]byte(dashboardPropsRaw), &dashboardProperties); err != nil {
		return fmt.Errorf("parsing JSON: %+v", err)
	}

	props.Properties = &dashboardProperties

	if _, err := client.CreateOrUpdate(ctx, id, props); err != nil {
		return fmt.Errorf("creating/updating %s %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())

		// `CreateOrUpdate` isn't a Long Running Operation and the Dashboard can briefly 404 after it returns, which
		// surfaces as `Provider produced inconsistent result after apply`, so wait for it to become available.
		pollerType := custompollers.NewDashboardPoller(client, id)
		poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for %s to become available: %+v", id, err)
		}
	}

	return resourcePortalDashboardRead(d, meta)
}
