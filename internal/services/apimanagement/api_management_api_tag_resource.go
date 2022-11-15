package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementApiTag() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiTagCreate,
		Read:   resourceApiManagementApiTagRead,
		Delete: resourceApiManagementApiTagDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApiTagID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApiManagementApiTagCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiId, err := parse.ApiID(d.Get("api_id").(string))
	if err != nil {
		return err
	}

	tagName := d.Get("name").(string)
	tagId := parse.NewTagID(subscriptionId, apiId.ResourceGroup, apiId.ServiceName, tagName)
	if err != nil {
		return err
	}

	id := parse.NewApiTagID(subscriptionId, apiId.ResourceGroup, apiId.ServiceName, apiId.Name, tagId.Name)

	tagExists, err := client.Get(ctx, apiId.ResourceGroup, apiId.ServiceName, tagId.ID())
	if err != nil {
		if !utils.ResponseWasNotFound(tagExists.Response) {
			return fmt.Errorf("checking for presence of Tag %q: %s", id, err)
		}
	}

	tagAssignmentExist, err := client.GetByAPI(ctx, apiId.ResourceGroup, apiId.ServiceName, apiId.Name, tagId.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(tagAssignmentExist.Response) {
			return fmt.Errorf("checking for presence of Tag Assignment %q: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(tagAssignmentExist.Response) {
		return tf.ImportAsExistsError("azurerm_api_management_api_tag", id.ID())
	}

	if _, err := client.AssignToAPI(ctx, apiId.ResourceGroup, apiId.ServiceName, apiId.Name, tagId.Name); err != nil {
		return fmt.Errorf("assigning to Api %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiTagRead(d, meta)
}

func resourceApiManagementApiTagRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiTagID(d.Id())
	if err != nil {
		return err
	}

	apiId := parse.NewApiID(subscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName)
	tagId := parse.NewTagID(subscriptionId, id.ResourceGroup, id.ServiceName, id.TagName)

	resp, err := client.GetByAPI(ctx, id.ResourceGroup, id.ServiceName, apiId.Name, tagId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("api_id", apiId.ID())
	d.Set("name", tagId.Name)

	return nil
}

func resourceApiManagementApiTagDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiTagID(d.Id())
	if err != nil {
		return err
	}
	apiId := parse.NewApiID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName)
	tagId := parse.NewTagID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.TagName)

	if _, err = client.DetachFromAPI(ctx, id.ResourceGroup, id.ServiceName, apiId.Name, tagId.Name); err != nil {
		return fmt.Errorf("detaching api tag %q: %+v", id, err)
	}

	return nil
}
