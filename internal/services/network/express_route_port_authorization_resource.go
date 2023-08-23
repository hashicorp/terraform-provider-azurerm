// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceExpressRoutePortAuthorization() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRoutePortAuthorizationCreate,
		Read:   resourceExpressRoutePortAuthorizationRead,
		Delete: resourceExpressRoutePortAuthorizationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ExpressRoutePortAuthorizationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
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
	client := meta.(*clients.Client).Network.ExpressRoutePortAuthorizationsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := parse.NewExpressRoutePortAuthorizationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_port_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ExpressRoutePortName, id.AuthorizationName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_express_route_port_authorization", id.ID())
		}
	}

	properties := network.ExpressRoutePortAuthorization{
		ExpressRoutePortAuthorizationPropertiesFormat: &network.ExpressRoutePortAuthorizationPropertiesFormat{},
	}

	// can run only one create/update/delete operation of expressRoutePort at the same time
	portID := parse.NewExpressRoutePortID(id.SubscriptionId, id.ResourceGroup, id.ExpressRoutePortName)
	locks.ByID(portID.ID())
	defer locks.UnlockByID(portID.ID())

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRoutePortName, id.AuthorizationName, properties)
	if err != nil {
		return fmt.Errorf("Creating/Updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for  %s to finish creating/updating: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRoutePortAuthorizationRead(d, meta)
}

func resourceExpressRoutePortAuthorizationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortAuthorizationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRoutePortAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRoutePortName, id.AuthorizationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AuthorizationName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("express_route_port_name", id.ExpressRoutePortName)

	if props := resp.ExpressRoutePortAuthorizationPropertiesFormat; props != nil {
		d.Set("authorization_key", props.AuthorizationKey)
		d.Set("authorization_use_status", string(props.AuthorizationUseStatus))
	}

	return nil
}

func resourceExpressRoutePortAuthorizationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortAuthorizationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRoutePortAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	portID := parse.NewExpressRoutePortID(id.SubscriptionId, id.ResourceGroup, id.ExpressRoutePortName)
	locks.ByID(portID.ID())
	defer locks.UnlockByID(portID.ID())

	future, err := client.Delete(ctx, id.ResourceGroup, id.ExpressRoutePortName, id.AuthorizationName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}
