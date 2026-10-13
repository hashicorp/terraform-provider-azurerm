// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_named_value

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementNamedValueRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namedvalue.ParseNamedValueID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NamedValueId)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)
			d.Set("secret", pointer.From(props.Secret))
			// API will not return `value` when `secret` is `true`, in which case we shall not set the `value`. Refer to the issue : #6688
			if props.Secret != nil && !*props.Secret {
				d.Set("value", pointer.From(props.Value))
			}
			if err := d.Set("value_from_key_vault", flattenApiManagementNamedValueKeyVault(props.KeyVault)); err != nil {
				return fmt.Errorf("setting `value_from_key_vault`: %+v", err)
			}
			d.Set("tags", pointer.From(props.Tags))
		}
	}

	return nil
}

func flattenApiManagementNamedValueKeyVault(input *namedvalue.KeyVaultContractProperties) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"secret_id":          pointer.From(input.SecretIdentifier),
			"identity_client_id": pointer.From(input.IdentityClientId),
		},
	}
}
