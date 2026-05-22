// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/restorables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceCosmosDbRestorableDatabaseAccounts() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCosmosDbRestorableDatabaseAccountsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"location": commonschema.LocationWithoutForceNew(),

			"accounts": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"api_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"creation_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"deletion_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"restorable_locations": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"creation_time": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"deletion_time": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"location": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"regional_database_account_instance_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCosmosDbRestorableDatabaseAccountsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.RestorablesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	locationName := d.Get("location").(string)
	id := restorables.NewRestorableDatabaseAccountID(meta.(*clients.Client).Account.SubscriptionId, locationName, "read")

	locationID := restorables.NewLocationID(id.SubscriptionId, location.Normalize(locationName))
	resp, err := client.RestorableDatabaseAccountsListByLocation(ctx, locationID)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("location", locationName)

	if resp.Model != nil {
		if v := resp.Model.Value; v != nil {
			if err := d.Set("accounts", flattenCosmosDbRestorableDatabaseAccounts(v, d.Get("name").(string))); err != nil {
				return fmt.Errorf("flattening `accounts`: %+v", err)
			}
		}
	}

	d.SetId(id.ID())

	return nil
}

func flattenCosmosDbRestorableDatabaseAccounts(input *[]restorables.RestorableDatabaseAccountGetResult, accountName string) []interface{} {
	result := make([]interface{}, 0)

	if input == nil {
		return result
	}

	for _, item := range *input {
		if props := item.Properties; props != nil && pointer.From(props.AccountName) == accountName {
			result = append(result, map[string]interface{}{
				"id":                   pointer.From(item.Id),
				"api_type":             pointer.From(props.ApiType),
				"creation_time":        pointer.From(props.CreationTime),
				"deletion_time":        pointer.From(props.DeletionTime),
				"restorable_locations": flattenCosmosDbRestorableDatabaseAccountsRestorableLocations(props.RestorableLocations),
			})
		}
	}

	return result
}

func flattenCosmosDbRestorableDatabaseAccountsRestorableLocations(input *[]restorables.RestorableLocationResource) []interface{} {
	result := make([]interface{}, 0)

	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, map[string]interface{}{
			"creation_time":                         pointer.From(item.CreationTime),
			"deletion_time":                         pointer.From(item.DeletionTime),
			"location":                              pointer.From(item.LocationName),
			"regional_database_account_instance_id": pointer.From(item.RegionalDatabaseAccountInstanceId),
		})
	}

	return result
}
