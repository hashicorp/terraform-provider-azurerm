// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementNamedValue() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementNamedValueCreateUpdate,
		Read:   resourceApiManagementNamedValueRead,
		Update: resourceApiManagementNamedValueCreateUpdate,
		Delete: resourceApiManagementNamedValueDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := namedvalue.ParseNamedValueID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"value_from_key_vault": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"value", "value_from_key_vault"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"secret_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},
						"identity_client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"value": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"value", "value_from_key_vault"},
			},

			"secret": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceApiManagementNamedValueCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := namedvalue.NewNamedValueID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
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

	parameters := namedvalue.NamedValueCreateContract{
		Properties: &namedvalue.NamedValueCreateContractProperties{
			DisplayName: d.Get("display_name").(string),
			Secret:      pointer.To(d.Get("secret").(bool)),
			KeyVault:    expandApiManagementNamedValueKeyVault(d.Get("value_from_key_vault").([]interface{})),
		},
	}

	if v, ok := d.GetOk("value"); ok {
		parameters.Properties.Value = pointer.To(v.(string))
	}

	if tags, ok := d.GetOk("tags"); ok {
		parameters.Properties.Tags = utils.ExpandStringSlice(tags.([]interface{}))
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, parameters, namedvalue.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementNamedValueRead(d, meta)
}

func resourceApiManagementNamedValueRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementNamedValueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namedvalue.ParseNamedValueID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, namedvalue.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf(" deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandApiManagementNamedValueKeyVault(inputs []interface{}) *namedvalue.KeyVaultContractCreateProperties {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	result := namedvalue.KeyVaultContractCreateProperties{
		SecretIdentifier: pointer.To(input["secret_id"].(string)),
	}

	if v := input["identity_client_id"].(string); v != "" {
		result.IdentityClientId = pointer.To(v)
	}

	return &result
}

func flattenApiManagementNamedValueKeyVault(input *namedvalue.KeyVaultContractProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"secret_id":          pointer.From(input.SecretIdentifier),
			"identity_client_id": pointer.From(input.IdentityClientId),
		},
	}
}
