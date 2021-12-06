package datalake

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakestore/2016-11-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
	subscriptionId := meta.(*clients.Client).Datalake.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if properties := model.Properties; properties != nil {
			tier := ""
			if properties.CurrentTier != nil {
				tier = string(*properties.CurrentTier)
			}
			d.Set("tier", tier)

			encryptionState := ""
			if properties.EncryptionState != nil {
				encryptionState = string(*properties.EncryptionState)
			}
			d.Set("encryption_state", encryptionState)

			firewallState := ""
			if properties.FirewallState != nil {
				firewallState = string(*properties.FirewallState)
			}
			d.Set("firewall_state", firewallState)

			firewallAllowAzureIps := ""
			if properties.FirewallAllowAzureIps != nil {
				firewallAllowAzureIps = string(*properties.FirewallAllowAzureIps)
			}
			d.Set("firewall_allow_azure_ips", firewallAllowAzureIps)

			if config := properties.EncryptionConfig; config != nil {
				d.Set("encryption_type", string(config.Type))
			}
		}

		return tags.FlattenAndSet(d, tagsHelper.Flatten(model.Tags))
	}
	return nil
}
