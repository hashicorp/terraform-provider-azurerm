// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apioperationtag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/tag"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationTag() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiOperationTagCreateUpdate,
		Read:   resourceApiManagementApiOperationTagRead,
		Update: resourceApiManagementApiOperationTagCreateUpdate,
		Delete: resourceApiManagementApiOperationTagDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apioperationtag.ParseOperationTagID(id)
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
	tagClient := meta.(*clients.Client).ApiManagement.TagClient
	client := meta.(*clients.Client).ApiManagement.ApiOperationTagClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiOperationId, err := apioperationtag.ParseOperationID(d.Get("api_operation_id").(string))
	if err != nil {
		return err
	}

	apiName := getApiName(apiOperationId.ApiId)

	id := apioperationtag.NewOperationTagID(subscriptionId, apiOperationId.ResourceGroupName, apiOperationId.ServiceName, apiName, apiOperationId.OperationId, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.TagGetByOperation(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Tag %q: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_operation_tag", id.ID())
		}
	}

	parameters := tag.TagCreateUpdateParameters{
		Properties: &tag.TagContractProperties{
			DisplayName: d.Get("display_name").(string),
		},
	}

	tagId := tag.NewTagID(subscriptionId, apiOperationId.ResourceGroupName, apiOperationId.ServiceName, d.Get("name").(string))
	if _, err := tagClient.CreateOrUpdate(ctx, tagId, parameters, tag.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if _, err := client.TagAssignToOperation(ctx, id); err != nil {
		return fmt.Errorf("assigning to operation %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiOperationTagRead(d, meta)
}

func resourceApiManagementApiOperationTagRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.ApiOperationTagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperationtag.ParseOperationTagID(d.Id())
	if err != nil {
		return err
	}

	apiName := getApiName(id.ApiId)

	newId := apioperationtag.NewOperationTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.OperationId, id.TagId)
	resp, err := client.TagGetByOperation(ctx, newId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", newId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", newId, err)
	}

	d.Set("api_operation_id", apioperationtag.NewOperationID(subscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.OperationId).ID())
	d.Set("name", id.TagId)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)
		}
	}

	return nil
}

func resourceApiManagementApiOperationTagDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationTagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperationtag.ParseOperationTagID(d.Id())
	if err != nil {
		return err
	}

	apiName := getApiName(id.ApiId)

	newId := apioperationtag.NewOperationTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.OperationId, id.TagId)
	if _, err = client.TagDetachFromOperation(ctx, newId); err != nil {
		return fmt.Errorf("deleting %q: %+v", newId, err)
	}

	return nil
}
