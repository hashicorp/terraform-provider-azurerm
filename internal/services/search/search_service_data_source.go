// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package search

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2024-06-01-preview/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2024-06-01-preview/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2024-06-01-preview/services"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"customer_managed_key_encryption_compliance_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

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

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"tags": commonschema.TagsDataSource(),
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
			partitionCount := 1
			replicaCount := 1
			publicNetworkAccess := true

			if props.EncryptionWithCmk != nil {
				d.Set("customer_managed_key_encryption_compliance_status", string(pointer.From(props.EncryptionWithCmk.EncryptionComplianceStatus)))
			}

			if count := props.PartitionCount; count != nil {
				partitionCount = int(*count)
			}

			if count := props.ReplicaCount; count != nil {
				replicaCount = int(*count)
			}

			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = strings.EqualFold(string(pointer.From(props.PublicNetworkAccess)), string(services.PublicNetworkAccessEnabled))
			}

			d.Set("partition_count", partitionCount)
			d.Set("replica_count", replicaCount)
			d.Set("public_network_access_enabled", publicNetworkAccess)
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err = d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	primaryKey := ""
	secondaryKey := ""
	adminKeysClient := meta.(*clients.Client).Search.AdminKeysClient
	adminKeysId, err := adminkeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}
	adminKeysResp, err := adminKeysClient.Get(ctx, *adminKeysId, adminkeys.GetOperationOptions{})
	if err != nil && !response.WasStatusCode(adminKeysResp.HttpResponse, http.StatusForbidden) {
		return fmt.Errorf("retrieving Admin Keys for %s: %+v", id, err)
	}
	if model := adminKeysResp.Model; model != nil {
		primaryKey = pointer.From(model.PrimaryKey)
		secondaryKey = pointer.From(model.SecondaryKey)
	}
	d.Set("primary_key", primaryKey)
	d.Set("secondary_key", secondaryKey)

	queryKeysClient := meta.(*clients.Client).Search.QueryKeysClient
	queryKeysId, err := querykeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}
	queryKeysResp, err := queryKeysClient.ListBySearchService(ctx, *queryKeysId, querykeys.ListBySearchServiceOperationOptions{})
	if err != nil && !response.WasStatusCode(queryKeysResp.HttpResponse, http.StatusForbidden) {
		return fmt.Errorf("retrieving Query Keys for %s: %+v", id, err)
	}
	if err := d.Set("query_keys", flattenSearchQueryKeys(queryKeysResp.Model)); err != nil {
		return fmt.Errorf("setting `query_keys`: %+v", err)
	}

	return nil
}
