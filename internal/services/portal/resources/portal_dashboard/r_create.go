// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_dashboard

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	portalClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePortalDashboardCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := portalClient.From(meta).DashboardsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := dashboard.NewDashboardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	return resourcePortalDashboardRead(d, meta)
}
