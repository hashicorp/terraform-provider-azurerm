// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/authorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwareExpressRouteAuthorization() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVmwareExpressRouteAuthorizationCreate,
		Read:   resourceVmwareExpressRouteAuthorizationRead,
		Delete: resourceVmwareExpressRouteAuthorizationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := authorizations.ParseAuthorizationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"private_cloud_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateCloudID,
			},

			"express_route_authorization_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"express_route_authorization_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceVmwareExpressRouteAuthorizationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.AuthorizationClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	privateCloudId, err := privateclouds.ParsePrivateCloudID(d.Get("private_cloud_id").(string))
	if err != nil {
		return err
	}

	id := authorizations.NewAuthorizationID(subscriptionId, privateCloudId.ResourceGroupName, privateCloudId.PrivateCloudName, name)
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_vmware_express_route_authorization", id.ID())
	}

	props := authorizations.ExpressRouteAuthorization{}

	if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVmwareExpressRouteAuthorizationRead(d, meta)
}

func resourceVmwareExpressRouteAuthorizationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.AuthorizationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := authorizations.ParseAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", *id, err)
	}

	d.Set("name", id.AuthorizationName)
	d.Set("private_cloud_id", privateclouds.NewPrivateCloudID(subscriptionId, id.ResourceGroupName, id.PrivateCloudName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("express_route_authorization_id", props.ExpressRouteAuthorizationId)
			d.Set("express_route_authorization_key", props.ExpressRouteAuthorizationKey)
		}
	}

	return nil
}

func resourceVmwareExpressRouteAuthorizationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.AuthorizationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := authorizations.ParseAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
