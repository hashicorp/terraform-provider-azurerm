// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name api_management_api_tag -service-package-name apimanagement -properties "name" -compare-values "resource_group_name:api_id,service_name:api_id,api_id:api_id" -known-values "subscription_id:data.Subscriptions.Primary" -test-name basic

func resourceApiManagementApiTag() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiTagCreate,
		Read:   resourceApiManagementApiTagRead,
		Delete: resourceApiManagementApiTagDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&apitag.ApiTagId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&apitag.ApiTagId{}),
		},

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
				ValidateFunc: api.ValidateApiID,
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
	tagClient := meta.(*clients.Client).ApiManagement.TagClient
	client := meta.(*clients.Client).ApiManagement.ApiTagClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiId, err := api.ParseApiID(d.Get("api_id").(string))
	if err != nil {
		return err
	}

	tagId := tag.NewTagID(subscriptionId, apiId.ResourceGroupName, apiId.ServiceName, d.Get("name").(string))

	id := apitag.NewApiTagID(subscriptionId, apiId.ResourceGroupName, apiId.ServiceName, apiId.ApiId, d.Get("name").(string))

	tagExists, err := tagClient.Get(ctx, tagId)
	if err != nil {
		if !response.WasNotFound(tagExists.HttpResponse) {
			return fmt.Errorf("checking for presence of Tag %q: %s", id, err)
		}
	}

	tagAssignmentExist, err := client.TagGetByApi(ctx, id)
	if err != nil {
		if !response.WasNotFound(tagAssignmentExist.HttpResponse) {
			return fmt.Errorf("checking for presence of Tag Assignment %q: %s", id, err)
		}
	}

	if !response.WasNotFound(tagAssignmentExist.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_api_management_api_tag", id.ID())
	}

	if _, err := client.TagAssignToApi(ctx, id); err != nil {
		return fmt.Errorf("assigning to Api %q: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceApiManagementApiTagRead(d, meta)
}

func resourceApiManagementApiTagRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.ApiTagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitag.ParseApiTagID(d.Id())
	if err != nil {
		return err
	}

	apiId := api.NewApiID(subscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId)
	tagId := apitag.NewApiTagID(subscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.TagId)

	resp, err := client.TagGetByApi(ctx, tagId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", tagId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", tagId, err)
	}

	d.Set("api_id", apiId.ID())
	d.Set("name", id.TagId)

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceApiManagementApiTagDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitag.ParseApiTagID(d.Id())
	if err != nil {
		return err
	}

	newId := apitag.NewApiTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.TagId)

	if _, err = client.TagDetachFromApi(ctx, newId); err != nil {
		return fmt.Errorf("detaching api tag %q: %+v", newId, err)
	}

	return nil
}
