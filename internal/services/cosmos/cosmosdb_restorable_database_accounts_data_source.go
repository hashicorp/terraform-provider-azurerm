// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
	client := meta.(*clients.Client).Cosmos.RestorableDatabaseAccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewRestorableDatabaseAccountID(subscriptionId, d.Get("location").(string), "read")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resp, err := client.ListByLocation(ctx, location)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("location", location)

	if props := resp.Value; props != nil {
		if err := d.Set("accounts", flattenCosmosDbRestorableDatabaseAccounts(props, name)); err != nil {
			return fmt.Errorf("flattening `accounts`: %+v", err)
		}
	}

	d.SetId(id.ID())

	return nil
}

func flattenCosmosDbRestorableDatabaseAccounts(input *[]documentdb.RestorableDatabaseAccountGetResult, accountName string) []interface{} {
	result := make([]interface{}, 0)

	if len(*input) == 0 {
		return result
	}

	for _, item := range *input {
		if props := item.RestorableDatabaseAccountProperties; props != nil && props.AccountName != nil && *props.AccountName == accountName {
			var id, creationTime, deletionTime string
			var apiType documentdb.APIType

			if item.ID != nil {
				id = *item.ID
			}

			if props.APIType != "" {
				apiType = props.APIType
			}

			if props.CreationTime != nil {
				creationTime = props.CreationTime.Format(time.RFC3339)
			}

			if props.DeletionTime != nil {
				deletionTime = props.DeletionTime.Format(time.RFC3339)
			}

			result = append(result, map[string]interface{}{
				"id":                   id,
				"api_type":             string(apiType),
				"creation_time":        creationTime,
				"deletion_time":        deletionTime,
				"restorable_locations": flattenCosmosDbRestorableDatabaseAccountsRestorableLocations(props.RestorableLocations),
			})
		}
	}

	return result
}

func flattenCosmosDbRestorableDatabaseAccountsRestorableLocations(input *[]documentdb.RestorableLocationResource) []interface{} {
	result := make([]interface{}, 0)

	if len(*input) == 0 {
		return result
	}

	for _, item := range *input {
		var location, regionalDatabaseAccountInstanceId, creationTime, deletionTime string

		if item.LocationName != nil {
			location = *item.LocationName
		}

		if item.RegionalDatabaseAccountInstanceID != nil {
			regionalDatabaseAccountInstanceId = *item.RegionalDatabaseAccountInstanceID
		}

		if item.CreationTime != nil {
			creationTime = item.CreationTime.Format(time.RFC3339)
		}

		if item.DeletionTime != nil {
			deletionTime = item.DeletionTime.Format(time.RFC3339)
		}

		result = append(result, map[string]interface{}{
			"creation_time":                         creationTime,
			"deletion_time":                         deletionTime,
			"location":                              location,
			"regional_database_account_instance_id": regionalDatabaseAccountInstanceId,
		})
	}

	return result
}
