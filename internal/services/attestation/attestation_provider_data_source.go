package attestation

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/sdk/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceAttestationProvider() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmAttestationProviderRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"attestation_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"trust_model": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmAttestationProviderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := attestationproviders.NewAttestationProvidersID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(parse.NewProviderID(subscriptionId, resourceGroup, name).ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if resp.Model != nil {
		d.Set("location", location.Normalize(resp.Model.Location))

		if props := resp.Model.Properties; props != nil {
			d.Set("attestation_uri", props.AttestUri)
			d.Set("trust_model", props.TrustModel)
		}
		return tags.FlattenAndSet(d, flattenTags(resp.Model.Tags))
	}

	return nil
}
