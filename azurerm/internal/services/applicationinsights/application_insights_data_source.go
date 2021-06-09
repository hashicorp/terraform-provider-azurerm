package applicationinsights

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApplicationInsights() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmApplicationInsightsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"instrumentation_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"application_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"app_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"retention_in_days": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceArmApplicationInsightsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Application Insights %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("retrieving Application Insights %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(parse.NewComponentID(subscriptionId, resGroup, name).ID())
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.ApplicationInsightsComponentProperties; props != nil {
		d.Set("app_id", props.AppID)
		d.Set("application_type", props.ApplicationType)
		d.Set("connection_string", props.ConnectionString)
		d.Set("instrumentation_key", props.InstrumentationKey)
		retentionInDays := 0
		if props.RetentionInDays != nil {
			retentionInDays = int(*props.RetentionInDays)
		}
		d.Set("retention_in_days", retentionInDays)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
