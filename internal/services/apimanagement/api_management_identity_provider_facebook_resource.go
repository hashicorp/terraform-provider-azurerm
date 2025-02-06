// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementIdentityProviderFacebook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementIdentityProviderFacebookCreateUpdate,
		Read:   resourceApiManagementIdentityProviderFacebookRead,
		Update: resourceApiManagementIdentityProviderFacebookCreateUpdate,
		Delete: resourceApiManagementIdentityProviderFacebookDelete,

		Importer: identityProviderImportFunc(identityprovider.IdentityProviderTypeFacebook),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"app_secret": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementIdentityProviderFacebookCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := identityprovider.NewIdentityProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), identityprovider.IdentityProviderTypeFacebook)
	clientID := d.Get("app_id").(string)
	clientSecret := d.Get("app_secret").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_identity_provider_facebook", id.ID())
		}
	}

	parameters := identityprovider.IdentityProviderCreateContract{
		Properties: &identityprovider.IdentityProviderCreateContractProperties{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			Type:         pointer.To(identityprovider.IdentityProviderTypeFacebook),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, identityprovider.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementIdentityProviderFacebookRead(d, meta)
}

func resourceApiManagementIdentityProviderFacebookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := identityprovider.ParseIdentityProviderID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroupName
	serviceName := id.ServiceName

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("app_id", props.ClientId)
		}
	}

	return nil
}

func resourceApiManagementIdentityProviderFacebookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := identityprovider.ParseIdentityProviderID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, identityprovider.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
