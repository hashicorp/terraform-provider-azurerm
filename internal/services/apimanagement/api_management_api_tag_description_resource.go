// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementApiTagDescription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiTagDescriptionCreateUpdate,
		Read:   resourceApiManagementApiTagDescriptionRead,
		Update: resourceApiManagementApiTagDescriptionCreateUpdate,
		Delete: resourceApiManagementApiTagDescriptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApiTagDescriptionsID(id)
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
				ValidateFunc: validate.ApiTagID,
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

	apiTagId, err := parse.ApiTagID(d.Get("api_tag_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_id`: %v", err)
	}

	id := parse.NewApiTagDescriptionsID(apiTagId.SubscriptionId, apiTagId.ResourceGroup, apiTagId.ServiceName, apiTagId.ApiName, apiTagId.TagName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_api_tag_description", id.ID())
		}
	}

	tagDescParameter := apimanagement.TagDescriptionCreateParameters{TagDescriptionBaseProperties: &apimanagement.TagDescriptionBaseProperties{}}
	if v, ok := d.GetOk("description"); ok {
		tagDescParameter.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_url"); ok {
		tagDescParameter.ExternalDocsURL = utils.String(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_description"); ok {
		tagDescParameter.ExternalDocsDescription = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName, tagDescParameter, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiTagDescriptionRead(d, meta)
}

func resourceApiManagementApiTagDescriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiTagDescriptionsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	apiTagId := parse.NewApiTagID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName)

	d.Set("api_tag_id", apiTagId.ID())
	d.Set("description", resp.Description)
	d.Set("external_documentation_url", resp.ExternalDocsURL)
	d.Set("external_documentation_description", resp.ExternalDocsDescription)

	return nil
}

func resourceApiManagementApiTagDescriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiTagDescriptionsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
