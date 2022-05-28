package portal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/portal/mgmt/2019-01-01-preview/portal"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourcePortalDashboardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	displayName, displayNameOk := d.GetOk("display_name")
	resourceGroup := d.Get("resource_group_name").(string)

	var dashboard portal.Dashboard

	if !displayNameOk {
		var err error
		dashboard, err = client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(dashboard.Response) {
				return fmt.Errorf("portal Dashboard %q was not found in Resource Group %q", name, resourceGroup)
			}
			return fmt.Errorf("retrieving Portal Dashboard %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	} else {
		dashboards := make([]portal.Dashboard, 0)

		iterator, err := client.ListByResourceGroupComplete(ctx, resourceGroup)
		if err != nil {
			if utils.ResponseWasNotFound(iterator.Response().Response) {
				return fmt.Errorf("no Portal Dashboards were found for Resource Group %q", resourceGroup)
			}
			return fmt.Errorf("getting list of Portal Dashboards (Resource Group %q): %+v", resourceGroup, err)
		}

		for iterator.NotDone() {
			dashboard = iterator.Value()

			found := false
			for k, v := range dashboard.Tags {
				if k == "hidden-title" && *v == displayName {
					found = true
					break
				}
			}

			if found {
				dashboards = append(dashboards, dashboard)
			}
			if err := iterator.NextWithContext(ctx); err != nil {
				return err
			}
		}

		if 1 > len(dashboards) {
			return fmt.Errorf("no Portal Dashboards were found for Resource Group %q", resourceGroup)
		}

		if len(dashboards) > 1 {
			return fmt.Errorf("multiple Portal Dashboards were found for Resource Group %q", resourceGroup)
		}

		dashboard = dashboards[0]
	}

	if dashboard.Name == nil {
		return fmt.Errorf("portal Dashboard name is empty in Resource Group %s", resourceGroup)
	}

	id := parse.NewDashboardID(subscriptionId, resourceGroup, *dashboard.Name)

	d.SetId(id.ID())

	d.Set("name", name)
	d.Set("display_name", displayName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(dashboard.Location))

	props, jsonErr := json.Marshal(dashboard.DashboardProperties)
	if jsonErr != nil {
		return fmt.Errorf("parsing JSON for Portal Dashboard Properties: %+v", jsonErr)
	}
	d.Set("dashboard_properties", string(props))

	return tags.FlattenAndSet(d, dashboard.Tags)
}
