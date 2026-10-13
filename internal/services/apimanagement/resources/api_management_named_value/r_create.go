// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_named_value

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementNamedValueCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := namedvalue.NewNamedValueID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf(" checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_property", id.ID())
			}
		}
	}

	parameters := namedvalue.NamedValueCreateContract{
		Properties: &namedvalue.NamedValueCreateContractProperties{
			DisplayName: d.Get("display_name").(string),
			Secret:      pointer.To(d.Get("secret").(bool)),
			KeyVault:    expandApiManagementNamedValueKeyVault(d.Get("value_from_key_vault").([]any)),
		},
	}

	if parameters.Properties.KeyVault != nil && (parameters.Properties.Secret == nil || !*parameters.Properties.Secret) {
		return errors.New("`secret` must be true when `value_from_key_vault` is set")
	}

	if v, ok := d.GetOk("value"); ok {
		parameters.Properties.Value = pointer.To(v.(string))
	}

	if tags, ok := d.GetOk("tags"); ok {
		parameters.Properties.Tags = pluginsdk.ExpandStringSlice(tags.([]any))
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, namedvalue.CreateOrUpdateOperationOptions{}, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, namedvalue.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceApiManagementNamedValueRead(d, meta)
}

func expandApiManagementNamedValueKeyVault(inputs []any) *namedvalue.KeyVaultContractCreateProperties {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]any)

	result := namedvalue.KeyVaultContractCreateProperties{
		SecretIdentifier: pointer.To(input["secret_id"].(string)),
	}

	if v := input["identity_client_id"].(string); v != "" {
		result.IdentityClientId = pointer.To(v)
	}

	return &result
}
