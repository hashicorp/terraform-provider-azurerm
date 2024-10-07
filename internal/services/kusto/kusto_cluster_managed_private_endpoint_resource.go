// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/managedprivateendpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoClusterManagedPrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterManagedPrivateEndpointCreateUpdate,
		Read:   resourceKustoClusterManagedPrivateEndpointRead,
		Update: resourceKustoClusterManagedPrivateEndpointCreateUpdate,
		Delete: resourceKustoClusterManagedPrivateEndpointDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoManagedPrivateEndpointV0ToV1{},
			1: migration.KustoManagedPrivateEndpointV1ToV2{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managedprivateendpoints.ParseManagedPrivateEndpointID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"private_link_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"private_link_resource_region": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"request_message": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceKustoClusterManagedPrivateEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := managedprivateendpoints.NewManagedPrivateEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		managedPrivateEndpoint, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(managedPrivateEndpoint.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(managedPrivateEndpoint.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_managed_private_endpoint", id.ID())
		}
	}

	managedPrivateEndpoint := managedprivateendpoints.ManagedPrivateEndpoint{
		Properties: &managedprivateendpoints.ManagedPrivateEndpointProperties{
			PrivateLinkResourceId: d.Get("private_link_resource_id").(string),
			GroupId:               d.Get("group_id").(string),
		},
	}

	if v, ok := d.GetOk("private_link_resource_region"); ok {
		managedPrivateEndpoint.Properties.PrivateLinkResourceRegion = utils.String(v.(string))
	}

	if v, ok := d.GetOk("request_message"); ok {
		managedPrivateEndpoint.Properties.RequestMessage = utils.String(v.(string))
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, managedPrivateEndpoint)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoClusterManagedPrivateEndpointRead(d, meta)
}

func resourceKustoClusterManagedPrivateEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(d.Id())
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

	d.Set("name", id.ManagedPrivateEndpointName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if resp.Model != nil {
		props := resp.Model.Properties
		if props != nil {
			d.Set("private_link_resource_id", props.PrivateLinkResourceId)
			d.Set("group_id", props.GroupId)

			if props.PrivateLinkResourceRegion != nil {
				d.Set("private_link_resource_region", props.PrivateLinkResourceRegion)
			}

			if props.RequestMessage != nil {
				d.Set("request_message", props.RequestMessage)
			}
		}
	}

	return nil
}

func resourceKustoClusterManagedPrivateEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
