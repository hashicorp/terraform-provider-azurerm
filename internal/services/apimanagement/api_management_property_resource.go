package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementProperty() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementPropertyCreateUpdate,
		Read:   resourceApiManagementPropertyRead,
		Update: resourceApiManagementPropertyCreateUpdate,
		Delete: resourceApiManagementPropertyDelete,

		DeprecationMessage: "This resource has been superseded by `azurerm_api_management_named_value` to reflects changes in the API/SDK and will be removed in version 3.0 of the provider.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PropertyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiPropertyV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"value": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

func resourceApiManagementPropertyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPropertyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.NamedValueName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %s", id, err)
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
			Value:       utils.String(d.Get("value").(string)),
		},
	}

	if tags, ok := d.GetOk("tags"); ok {
		parameters.NamedValueCreateContractProperties.Tags = utils.ExpandStringSlice(tags.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.NamedValueName, parameters, "")
	if err != nil {
		return fmt.Errorf(" creating or updating Property %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementPropertyRead(d, meta)
}

func resourceApiManagementPropertyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PropertyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.NamedValueName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", *id, err)
	}

	d.Set("name", id.NamedValueName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	if properties := resp.NamedValueContractProperties; properties != nil {
		d.Set("display_name", properties.DisplayName)
		d.Set("secret", properties.Secret)
		if !*properties.Secret {
			d.Set("value", properties.Value)
		}
		d.Set("tags", properties.Tags)
	}

	return nil
}

func resourceApiManagementPropertyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.NamedValueClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PropertyID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.NamedValueName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
