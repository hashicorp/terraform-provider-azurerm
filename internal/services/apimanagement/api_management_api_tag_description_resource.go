// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitagdescription"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiTagDescription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiTagDescriptionCreateUpdate,
		Read:   resourceApiManagementApiTagDescriptionRead,
		Update: resourceApiManagementApiTagDescriptionCreateUpdate,
		Delete: resourceApiManagementApiTagDescriptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apitagdescription.ParseTagDescriptionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"api_tag_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: apitag.ValidateApiTagID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"external_documentation_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"external_documentation_description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementApiTagDescriptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiTagId, err := apitag.ParseApiTagID(d.Get("api_tag_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_id`: %v", err)
	}

	apiName := getApiName(apiTagId.ApiId)

	id := apitagdescription.NewTagDescriptionID(apiTagId.SubscriptionId, apiTagId.ResourceGroupName, apiTagId.ServiceName, apiName, apiTagId.TagId)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_tag_description", id.ID())
		}
	}

	tagDescParameter := apitagdescription.TagDescriptionCreateParameters{Properties: &apitagdescription.TagDescriptionBaseProperties{}}
	if v, ok := d.GetOk("description"); ok {
		tagDescParameter.Properties.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_url"); ok {
		tagDescParameter.Properties.ExternalDocsURL = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_description"); ok {
		tagDescParameter.Properties.ExternalDocsDescription = pointer.To(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, tagDescParameter, apitagdescription.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiTagDescriptionRead(d, meta)
}

func resourceApiManagementApiTagDescriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitagdescription.ParseTagDescriptionID(d.Id())
	if err != nil {
		return err
	}

	apiName := getApiName(id.ApiId)

	newId := apitagdescription.NewTagDescriptionID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.TagDescriptionId)
	resp, err := client.Get(ctx, newId)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", newId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", newId, err)
	}

	apiTagId := apitag.NewApiTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.TagDescriptionId)

	d.Set("api_tag_id", apiTagId.ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("external_documentation_url", pointer.From(props.ExternalDocsURL))
			d.Set("external_documentation_description", pointer.From(props.ExternalDocsDescription))
		}
	}

	return nil
}

func resourceApiManagementApiTagDescriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitagdescription.ParseTagDescriptionID(d.Id())
	if err != nil {
		return err
	}

	name := getApiName(id.ApiId)

	newId := apitagdescription.NewTagDescriptionID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, name, id.TagDescriptionId)
	resp, err := client.Delete(ctx, newId, apitagdescription.DeleteOperationOptions{})
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", newId, err)
		}
	}

	return nil
}
