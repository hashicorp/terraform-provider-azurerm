// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apirelease"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiRelease() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiReleaseCreateUpdate,
		Read:   resourceApiManagementApiReleaseRead,
		Update: resourceApiManagementApiReleaseCreateUpdate,
		Delete: resourceApiManagementApiReleaseDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apirelease.ParseReleaseID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementChildName,
			},

			"api_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: api.ValidateApiID,
			},

			"notes": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementApiReleaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	apiId, err := api.ParseApiID(d.Get("api_id").(string))
	if err != nil {
		return err
	}

	id := apirelease.NewReleaseID(subscriptionId, apiId.ResourceGroupName, apiId.ServiceName, apiId.ApiId, name)
	ifMatch := "*"

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_release", id.ID())
		}
		ifMatch = ""
	}

	parameters := apirelease.ApiReleaseContract{
		Properties: &apirelease.ApiReleaseContractProperties{
			ApiId: pointer.To(d.Get("api_id").(string)),
			Notes: pointer.To(d.Get("notes").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apirelease.CreateOrUpdateOperationOptions{IfMatch: pointer.To(ifMatch)}); err != nil {
		return fmt.Errorf("creating/ updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementApiReleaseRead(d, meta)
}

func resourceApiManagementApiReleaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apirelease.ParseReleaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] apimanagement %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.ReleaseId)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("api_id", api.NewApiID(subscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId).ID())
			d.Set("notes", pointer.From(props.Notes))
		}
	}
	return nil
}

func resourceApiManagementApiReleaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apirelease.ParseReleaseID(d.Id())
	if err != nil {
		return err
	}

	name := getApiName(id.ApiId)

	newId := apirelease.NewReleaseID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, name, id.ReleaseId)
	if _, err := client.Delete(ctx, newId, apirelease.DeleteOperationOptions{IfMatch: pointer.To("*")}); err != nil {
		return fmt.Errorf("deleting %s: %+v", newId, err)
	}
	return nil
}
