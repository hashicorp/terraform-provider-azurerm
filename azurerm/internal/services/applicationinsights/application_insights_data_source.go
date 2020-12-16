package applicationinsights

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApplicationInsights() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmApplicationInsightsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"instrumentation_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"application_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"retention_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmApplicationInsightsRead(d *schema.ResourceData, meta interface{}) error {
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
