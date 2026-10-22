// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maps_account

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMapsAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_maps_account", id.ID())
		}
	}

	dataStores, err := expandDataStore(d.Get("data_store").([]any))
	if err != nil {
		return fmt.Errorf("expanding `data_store`: %+v", err)
	}

	parameters := accounts.MapsAccount{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: accounts.Sku{
			Name: accounts.Name(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &accounts.MapsAccountProperties{
			DisableLocalAuth: pointer.To(!d.Get("local_authentication_enabled").(bool)),
			Cors:             expandCors(d.Get("cors").([]any)),
			LinkedResources:  dataStores,
		},
	}

	// setting anything into identity returns a 400 Bad Request error if the location of the maps account is `global`
	if v, ok := d.GetOk("identity"); ok {
		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(v.([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = identityExpanded
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// TODO: add a Pandora workaround for this
	// These should actually be LROs, but they're not, custom poller is required until https://github.com/Azure/azure-rest-api-specs/issues/29501 is resolved
	pollerType := custompollers.NewMapsAccountPoller(client, id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return resourceMapsAccountRead(d, meta)
}

func expandCors(input []any) *accounts.CorsRules {
	if len(input) == 0 {
		return nil
	}

	cors := input[0].(map[string]any)

	corsRule := make([]accounts.CorsRule, 0)

	corsRule = append(corsRule, accounts.CorsRule{
		AllowedOrigins: pointer.From(pluginsdk.ExpandStringSlice(cors["allowed_origins"].([]any))),
	})

	return &accounts.CorsRules{
		CorsRules: &corsRule,
	}
}

func expandDataStore(input []any) (*[]accounts.LinkedResource, error) {
	if len(input) == 0 {
		return nil, nil
	}

	linkedResources := make([]accounts.LinkedResource, 0)

	for _, i := range input {
		dataStore := i.(map[string]any)

		storageAccountId, err := commonids.ParseStorageAccountID(dataStore["storage_account_id"].(string))
		if err != nil {
			return nil, err
		}

		linkedResources = append(linkedResources, accounts.LinkedResource{
			Id:         storageAccountId.ID(),
			UniqueName: dataStore["unique_name"].(string),
		})
	}

	return &linkedResources, nil
}
