package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCosmosDbRestorableDatabaseAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCosmosDbRestorableDatabaseAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"location": location.SchemaWithoutForceNew(),

			"restorable_db_account_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceCosmosDbRestorableDatabaseAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		restorableDbAccountIds := make([]string, 0)
		for _, v := range *props {
			if v.ID != nil && v.RestorableDatabaseAccountProperties != nil && v.RestorableDatabaseAccountProperties.AccountName != nil && *v.RestorableDatabaseAccountProperties.AccountName == name {
				restorableDbAccountIds = append(restorableDbAccountIds, *v.ID)
			}
		}
		d.Set("restorable_db_account_ids", restorableDbAccountIds)
	}

	d.SetId(id.ID())

	return nil
}
