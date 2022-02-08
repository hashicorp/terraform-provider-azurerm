package portal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
				Required:     true,
				ValidateFunc: validate.DashboardName,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
			"location":            azure.SchemaLocationForDataSource(),
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
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDashboardID(subscriptionId, resourceGroup, name)
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("portal dashboard %q was not found in Resource Group %q", id.Name, id.ResourceGroup)
		}
		return fmt.Errorf("retrieving Portal Dashboard %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	props, jsonErr := json.Marshal(resp.DashboardProperties)
	if jsonErr != nil {
		return fmt.Errorf("parsing JSON for Portal Dashboard Properties: %+v", jsonErr)
	}
	d.Set("dashboard_properties", string(props))

	return tags.FlattenAndSet(d, resp.Tags)
}
