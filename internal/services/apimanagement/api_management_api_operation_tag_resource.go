// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tag"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationTag() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceApiManagementApiOperationTagCreateUpdate,
		Read:   resourceApiManagementApiOperationTagRead,
		Delete: resourceApiManagementApiOperationTagDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apioperationtag.ParseOperationTagID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
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
		},
	}

	if !features.FivePointOhBeta() {
		resource.Update = resourceApiManagementApiOperationTagCreateUpdate
		resource.Timeouts.Update = pluginsdk.DefaultTimeout(30 * time.Minute)

		resource.Schema["display_name"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Computed:   true,
			Deprecated: "This property has been deprecated and will be removed in v5.0 of the provider. Use display_name property of azurerm_api_management_tag resource.",
		}
	}

	return resource
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
	tagId := tag.NewTagID(subscriptionId, apiOperationId.ResourceGroupName, apiOperationId.ServiceName, d.Get("name").(string))

	// For 4.0 we continue to create the tag if display_name is set (backward compatibility)
	displayName := d.Get("display_name").(string)
	if !features.FivePointOhBeta() && len(displayName) > 0 {
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

		if _, err := tagClient.CreateOrUpdate(ctx, tagId, parameters, tag.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf("creating/updating %q: %+v", id, err)
		}
	} else {
		tagAssignmentExist, err := client.TagGetByOperation(ctx, id)
		if err != nil {
			if !response.WasNotFound(tagAssignmentExist.HttpResponse) {
				return fmt.Errorf("checking for presence of Tag Assignment %q: %s", id, err)
			}
		}

		if !response.WasNotFound(tagAssignmentExist.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_operation_tag", id.ID())
		}

		tagExists, err := tagClient.Get(ctx, tagId)
		if err != nil {
			if !response.WasNotFound(tagExists.HttpResponse) {
				return fmt.Errorf("checking for presence of Tag %q: %s", id, err)
			}
		}
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

	if !features.FivePointOhBeta() {
		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				d.Set("display_name", props.DisplayName)
			}
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
