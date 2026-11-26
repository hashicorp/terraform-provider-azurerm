// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maps_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMapsAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku_name", string(model.Sku.Name))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		d.Set("identity", identityFlattened)

		if props := model.Properties; props != nil {
			d.Set("x_ms_client_id", props.UniqueId)
			d.Set("cors", flattenCors(props.Cors))

			dataStore, err := flattenDataStore(props.LinkedResources)
			if err != nil {
				return fmt.Errorf("flattening `data_store`: %+v", err)
			}
			d.Set("data_store", dataStore)

			localAuthenticationEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthenticationEnabled = !*props.DisableLocalAuth
			}
			d.Set("local_authentication_enabled", localAuthenticationEnabled)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	keysResp, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Access Keys for %s: %+v", *id, err)
	}
	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	return nil
}

func flattenCors(input *accounts.CorsRules) []any {
	output := make([]any, 0)

	if input == nil || input.CorsRules == nil || len(*input.CorsRules) == 0 {
		return output
	}

	// although this is a slice, only one element can be supplied/is present
	allowedOrigins := (*input.CorsRules)[0].AllowedOrigins

	output = append(output, map[string]any{
		"allowed_origins": allowedOrigins,
	})

	return output
}

func flattenDataStore(input *[]accounts.LinkedResource) ([]any, error) {
	output := make([]any, 0)

	if input == nil || len(*input) == 0 {
		return output, nil
	}

	for _, resource := range *input {
		storageAccountId, err := commonids.ParseStorageAccountID(resource.Id)
		if err != nil {
			return nil, err
		}

		output = append(output, map[string]any{
			"storage_account_id": storageAccountId.ID(),
			"unique_name":        resource.UniqueName,
		})
	}

	return output, nil
}
