package mixedreality

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mixedreality/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSpatialAnchorsAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpatialAnchorsAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w._()]+$`),
					"Spatial Anchors Account name must be 1 - 90 characters long, contain only word characters and underscores.",
				),
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSpatialAnchorsAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewSpatialAnchorsAccountID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("retrieving Spatial Anchors Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("account_domain", resp.AccountDomain)
	d.Set("account_id", resp.AccountID)

	return nil
}
