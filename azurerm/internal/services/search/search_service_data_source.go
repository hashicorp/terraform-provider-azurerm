package search

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/search/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSearchService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSearchServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"replica_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"partition_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"query_keys": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSearchServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSearchServiceID(subscriptionID, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Search Service %q (Resource Group %q) was not found", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("Error reading Search Service: %+v", err)
	}

	d.SetId(id.ID())

	if props := resp.ServiceProperties; props != nil {
		if count := props.PartitionCount; count != nil {
			d.Set("partition_count", int(*count))
		}

		if count := props.ReplicaCount; count != nil {
			d.Set("replica_count", int(*count))
		}

		d.Set("public_network_access_enabled", props.PublicNetworkAccess != "Disabled")
	}

	adminKeysClient := meta.(*clients.Client).Search.AdminKeysClient
	adminKeysResp, err := adminKeysClient.Get(ctx, id.ResourceGroup, id.Name, nil)
	if err == nil {
		d.Set("primary_key", adminKeysResp.PrimaryKey)
		d.Set("secondary_key", adminKeysResp.SecondaryKey)
	}

	queryKeysClient := meta.(*clients.Client).Search.QueryKeysClient
	queryKeysResp, err := queryKeysClient.ListBySearchService(ctx, id.ResourceGroup, id.Name, nil)
	if err == nil {
		d.Set("query_keys", flattenSearchQueryKeys(queryKeysResp.Values()))
	}

	if err := d.Set("identity", flattenSearchServiceIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	return nil
}
