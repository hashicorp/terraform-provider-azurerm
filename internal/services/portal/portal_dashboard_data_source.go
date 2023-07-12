// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package portal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourcePortalDashboard() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePortalDashboardRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.DashboardName,
				ExactlyOneOf: []string{"name", "display_name"},
			},
			"display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
			"location":            commonschema.LocationComputed(),
			"dashboard_properties": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},
			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourcePortalDashboardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	displayName := d.Get("display_name")

	id := dashboard.NewDashboardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	props := dashboard.Dashboard{}

	if displayName == "" {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("%s was not found", id)
			}
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if model := resp.Model; model != nil {
			props = *model
		}
	} else {
		dashboards := make([]dashboard.Dashboard, 0)

		resourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName)
		iterator, err := client.ListByResourceGroupComplete(ctx, resourceGroupId)
		if err != nil {
			return fmt.Errorf("getting list of Portal Dashboards within %s: %+v", resourceGroupId, err)
		}

		log.Printf("portal_debug iterator: %+v", iterator.Items)

		for _, item := range iterator.Items {
			tags := *item.Tags
			for k, v := range tags {
				if k == "hidden-title" && v == displayName {
					dashboards = append(dashboards, item)
					break
				}
			}
		}

		if len(dashboards) == 0 {
			return fmt.Errorf("no Portal Dashboards were found within %s", resourceGroupId)
		}

		if len(dashboards) > 1 {
			return fmt.Errorf("multiple (%d) Portal Dashboards were found within %s", len(dashboards), resourceGroupId)
		}

		props = dashboards[0]
		id.DashboardName = pointer.From(props.Name)
	}

	d.SetId(id.ID())

	d.Set("name", id.DashboardName)
	d.Set("display_name", displayName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.Normalize(props.Location))

	v, err := json.Marshal(props.Properties)
	if err != nil {
		return fmt.Errorf("parsing JSON for Portal Dashboard Properties: %+v", err)
	}
	d.Set("dashboard_properties", string(v))

	return tags.FlattenAndSet(d, props.Tags)
}
