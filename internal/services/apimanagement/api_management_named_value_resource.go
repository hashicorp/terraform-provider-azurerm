package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
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
			_, err := parse.NamedValueID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

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

	id := parse.NewNamedValueID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf(" checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_property", id.ID())
		}
	}

	parameters := apimanagement.NamedValueCreateContract{
		NamedValueCreateContractProperties: &apimanagement.NamedValueCreateContractProperties{
			DisplayName: utils.String(d.Get("display_name").(string)),
			Secret:      utils.Bool(d.Get("secret").(bool)),
			KeyVault:    expandApiManagementNamedValueKeyVault(d.Get("value_from_key_vault").([]interface{})),
		},
	}

	if v, ok := d.GetOk("value"); ok {
		parameters.NamedValueCreateContractProperties.Value = utils.String(v.(string))
	}

	if tags, ok := d.GetOk("tags"); ok {
		parameters.NamedValueCreateContractProperties.Tags = utils.ExpandStringSlice(tags.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.Name, parameters, "")
	if err != nil {
		return fmt.Errorf(" creating or updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementNamedValueRead(d, meta)
}

func resourceApiManagementNamedValueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamedValueID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf(" making Read request for %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	if properties := resp.NamedValueContractProperties; properties != nil {
		d.Set("display_name", properties.DisplayName)
		d.Set("secret", properties.Secret)
		// API will not return `value` when `secret` is `true`, in which case we shall not set the `value`. Refer to the issue : #6688
		if properties.Secret != nil && !*properties.Secret {
			d.Set("value", properties.Value)
		}
		if err := d.Set("value_from_key_vault", flattenApiManagementNamedValueKeyVault(properties.KeyVault)); err != nil {
			return fmt.Errorf("setting `value_from_key_vault`: %+v", err)
		}
		d.Set("tags", properties.Tags)
	}

	return nil
}

func resourceApiManagementNamedValueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamedValueID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf(" deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandApiManagementNamedValueKeyVault(inputs []interface{}) *apimanagement.KeyVaultContractCreateProperties {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	result := apimanagement.KeyVaultContractCreateProperties{
		SecretIdentifier: utils.String(input["secret_id"].(string)),
	}

	if v := input["identity_client_id"].(string); v != "" {
		result.IdentityClientID = utils.String(v)
	}

	return &result
}

func flattenApiManagementNamedValueKeyVault(input *apimanagement.KeyVaultContractProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var secretId, clientId string
	if input.SecretIdentifier != nil {
		secretId = *input.SecretIdentifier
	}

	if input.IdentityClientID != nil {
		clientId = *input.IdentityClientID
	}

	return []interface{}{
		map[string]interface{}{
			"secret_id":          secretId,
			"identity_client_id": clientId,
		},
	}
}
