package search

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/services"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"identity": commonschema.SystemAssignedIdentityComputed(),
		},
	}
}

func dataSourceSearchServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := services.NewSearchServiceID(subscriptionID, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id, services.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id.ID())
		}

		return fmt.Errorf("reading Search Service: %+v", err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			partitionCount := 0
			replicaCount := 0
			publicNetworkAccess := false

			if count := props.PartitionCount; count != nil {
				partitionCount = int(*count)
			}

			if count := props.ReplicaCount; count != nil {
				replicaCount = int(*count)
			}

			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = *props.PublicNetworkAccess != "Disabled"
			}

			d.Set("partition_count", partitionCount)
			d.Set("replica_count", replicaCount)
			d.Set("public_network_access_enabled", publicNetworkAccess)
		}

		if err = d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}
	}

	adminKeysClient := meta.(*clients.Client).Search.AdminKeysClient
	adminKeysId, err := adminkeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	adminKeysResp, err := adminKeysClient.Get(ctx, *adminKeysId, adminkeys.GetOperationOptions{})
	if err == nil {
		if model := adminKeysResp.Model; model != nil {
			d.Set("primary_key", model.PrimaryKey)
			d.Set("secondary_key", model.SecondaryKey)
		}
	}

	queryKeysClient := meta.(*clients.Client).Search.QueryKeysClient
	queryKeysId, err := querykeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}
	queryKeysResp, err := queryKeysClient.ListBySearchService(ctx, *queryKeysId, querykeys.ListBySearchServiceOperationOptions{})
	if err == nil {
		if model := queryKeysResp.Model; model != nil {
			d.Set("query_keys", flattenSearchQueryKeys(*model))

		}
	}

	return nil
}
