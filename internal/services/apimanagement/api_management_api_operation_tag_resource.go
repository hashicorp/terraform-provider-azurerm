// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tag"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationTag() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceApiManagementApiOperationTagCreate,
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
				ValidateFunc: validation.AsGeneratedID(apioperation.ParseOperationIDInsensitively),
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementChildName,
			},
		},
	}

	if !features.SixPointOh() {
		resource.Update = resourceApiManagementApiOperationTagUpdate
		resource.Timeouts.Update = pluginsdk.DefaultTimeout(30 * time.Minute)

		resource.Schema["display_name"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C when omitted the existing tag is assigned as is, so its display name is read back from Azure
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Deprecated:   "`display_name` has been deprecated in favour of the `azurerm_api_management_tag` resource and will be removed in v6.0 of the AzureRM Provider",
		}
	}

	return resource
}

func resourceApiManagementApiOperationTagCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	tagClient := meta.(*clients.Client).ApiManagement.TagClient
	client := meta.(*clients.Client).ApiManagement.ApiOperationTagClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiOperationId, err := apioperationtag.ParseOperationID(d.Get("api_operation_id").(string))
	if err != nil {
		return err
	}

	apiName := getApiName(apiOperationId.ApiId)

	id := apioperationtag.NewOperationTagID(subscriptionId, apiOperationId.ResourceGroupName, apiOperationId.ServiceName, apiName, apiOperationId.OperationId, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

	tagId := tag.NewTagID(subscriptionId, apiOperationId.ResourceGroupName, apiOperationId.ServiceName, d.Get("name").(string))

	if !features.SixPointOh() {
		if displayName := d.Get("display_name").(string); displayName != "" {
			if err := setApiManagementTagDisplayName(ctx, tagClient, tagId, displayName); err != nil {
				return err
			}
		}
	}

	existingTag, err := tagClient.Get(ctx, tagId)
	if err != nil {
		if response.WasNotFound(existingTag.HttpResponse) {
			return fmt.Errorf("%s was not found, create it with the `azurerm_api_management_tag` resource before assigning it to an operation", tagId)
		}

		return fmt.Errorf("retrieving %s: %+v", tagId, err)
	}

	if _, err := client.TagAssignToOperation(ctx, id); err != nil {
		return fmt.Errorf("assigning to operation %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiOperationTagRead(d, meta)
}

func resourceApiManagementApiOperationTagRead(d *pluginsdk.ResourceData, meta any) error {
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

	if !features.SixPointOh() {
		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				d.Set("display_name", props.DisplayName)
			}
		}
	}

	return nil
}

func resourceApiManagementApiOperationTagUpdate(d *pluginsdk.ResourceData, meta any) error {
	tagClient := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperationtag.ParseOperationTagID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("display_name") {
		tagId := tag.NewTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.TagId)
		if err := setApiManagementTagDisplayName(ctx, tagClient, tagId, d.Get("display_name").(string)); err != nil {
			return err
		}
	}

	return resourceApiManagementApiOperationTagRead(d, meta)
}

func resourceApiManagementApiOperationTagDelete(d *pluginsdk.ResourceData, meta any) error {
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

// setApiManagementTagDisplayName creates or updates the tag itself, which this resource only does for the deprecated `display_name`
func setApiManagementTagDisplayName(ctx context.Context, client *tag.TagClient, id tag.TagId, displayName string) error {
	parameters := tag.TagCreateUpdateParameters{
		Properties: &tag.TagContractProperties{
			DisplayName: displayName,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, tag.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	return nil
}
