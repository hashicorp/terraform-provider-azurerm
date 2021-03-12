package attestation

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/attestation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAttestationProvider() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAttestationProviderRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"attestation_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"trust_model": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmAttestationProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Attestation Provider %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("retrieving Attestation Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(parse.NewProviderID(subscriptionId, resourceGroup, name).ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.StatusResult; props != nil {
		d.Set("attestation_uri", props.AttestURI)
		d.Set("trust_model", props.TrustModel)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
