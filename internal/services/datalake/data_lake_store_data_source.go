package datalake

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceDataLakeStoreAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmDateLakeStoreAccountRead,

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

			"tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"firewall_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"firewall_allow_azure_ips": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDateLakeStoreAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	subscriptionId := meta.(*clients.Client).Datalake.StoreAccountsClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Data Lake Store Account %s was not found", id)
		}
		return fmt.Errorf("retrieving Data Lake Store Account %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.DataLakeStoreAccountProperties; properties != nil {
		d.Set("tier", string(properties.CurrentTier))

		d.Set("encryption_state", string(properties.EncryptionState))
		d.Set("firewall_allow_azure_ips", string(properties.FirewallAllowAzureIps))
		d.Set("firewall_state", string(properties.FirewallState))

		if config := properties.EncryptionConfig; config != nil {
			d.Set("encryption_type", string(config.Type))
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
