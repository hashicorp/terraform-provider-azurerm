package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiOperationTag() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiOperationTagCreateUpdate,
		Read:   resourceApiManagementApiOperationTagRead,
		Update: resourceApiManagementApiOperationTagCreateUpdate,
		Delete: resourceApiManagementApiOperationTagDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.OperationTagID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_operation_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiOperationID,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementChildName,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementApiOperationTagCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiOperationId, err := parse.ApiOperationID(d.Get("api_operation_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	id := parse.NewOperationTagID(subscriptionId, apiOperationId.ResourceGroup, apiOperationId.ServiceName, apiOperationId.ApiName, apiOperationId.OperationName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, apiOperationId.ResourceGroup, apiOperationId.ServiceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Tag %q: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_api_operation_tag", id.ID())
		}
	}

	parameters := apimanagement.TagCreateUpdateParameters{
		TagContractProperties: &apimanagement.TagContractProperties{
			DisplayName: utils.String(d.Get("display_name").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, apiOperationId.ResourceGroup, apiOperationId.ServiceName, name, parameters, ""); err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if _, err := client.AssignToOperation(ctx, apiOperationId.ResourceGroup, apiOperationId.ServiceName, apiOperationId.ApiName, apiOperationId.OperationName, name); err != nil {
		return fmt.Errorf("assigning to operation %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiOperationTagRead(d, meta)
}

func resourceApiManagementApiOperationTagRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.OperationTagID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.TagName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("api_operation_id", parse.NewApiOperationID(subscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName, id.OperationName).ID())
	d.Set("name", id.TagName)

	if props := resp.TagContractProperties; props != nil {
		d.Set("display_name", props.DisplayName)
	}

	return nil
}

func resourceApiManagementApiOperationTagDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.OperationTagID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.TagName, ""); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
