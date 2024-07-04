// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteportauthorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRoutePortAuthorization() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRoutePortAuthorizationCreate,
		Read:   resourceExpressRoutePortAuthorizationRead,
		Delete: resourceExpressRoutePortAuthorizationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := expressrouteportauthorizations.ParseExpressRoutePortAuthorizationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"express_route_port_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"authorization_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"authorization_use_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceExpressRoutePortAuthorizationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortAuthorizations
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := expressrouteportauthorizations.NewExpressRoutePortAuthorizationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_port_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_express_route_port_authorization", id.ID())
		}
	}

	properties := expressrouteportauthorizations.ExpressRoutePortAuthorization{
		Properties: &expressrouteportauthorizations.ExpressRoutePortAuthorizationPropertiesFormat{},
	}

	// can run only one create/update/delete operation of expressRoutePort at the same time
	portID := expressrouteports.NewExpressRoutePortID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRoutePortName)
	locks.ByID(portID.ID())
	defer locks.UnlockByID(portID.ID())

	if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRoutePortAuthorizationRead(d, meta)
}

func resourceExpressRoutePortAuthorizationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortAuthorizations
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteportauthorizations.ParseExpressRoutePortAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AuthorizationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("express_route_port_name", id.ExpressRoutePortName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("authorization_key", props.AuthorizationKey)
			d.Set("authorization_use_status", string(pointer.From(props.AuthorizationUseStatus)))
		}
	}

	return nil
}

func resourceExpressRoutePortAuthorizationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortAuthorizations
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteportauthorizations.ParseExpressRoutePortAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	portID := expressrouteports.NewExpressRoutePortID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRoutePortName)
	locks.ByID(portID.ID())
	defer locks.UnlockByID(portID.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
