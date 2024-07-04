// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorPrivateLinkScopedService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorPrivateLinkScopedServiceCreate,
		Read:   resourceMonitorPrivateLinkScopedServiceRead,
		Delete: resourceMonitorPrivateLinkScopedServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privatelinkscopedresources.ParseScopedResourceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"scope_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateLinkScopeName,
			},

			"linked_resource_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.Any(
					components.ValidateComponentID,
					workspaces.ValidateWorkspaceID,
					datacollectionendpoints.ValidateDataCollectionEndpointID,
				),
			},
		},
	}
}

func resourceMonitorPrivateLinkScopedServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Monitor.PrivateLinkScopedResourcesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privatelinkscopedresources.NewScopedResourceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("scope_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_private_link_scoped_service", id.ID())
		}
	}

	parameters := privatelinkscopedresources.ScopedResource{
		Properties: &privatelinkscopedresources.ScopedResourceProperties{
			LinkedResourceId: utils.String(d.Get("linked_resource_id").(string)),
		},
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorPrivateLinkScopedServiceRead(d, meta)
}

func resourceMonitorPrivateLinkScopedServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopedResourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopedresources.ParseScopedResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ScopedResourceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("scope_name", id.PrivateLinkScopeName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("linked_resource_id", normalizeLinkedResourceId(props.LinkedResourceId))
		}
	}

	return nil
}

func resourceMonitorPrivateLinkScopedServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopedResourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopedresources.ParseScopedResourceID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func normalizeLinkedResourceId(input *string) *string {
	if input == nil {
		return input
	}

	if resourceId, err := components.ParseComponentIDInsensitively(*input); err == nil {
		nomalizedId := resourceId.ID()
		return &nomalizedId
	}
	if resourceId, err := workspaces.ParseWorkspaceIDInsensitively(*input); err == nil {
		nomalizedId := resourceId.ID()
		return &nomalizedId
	}
	if resourceId, err := datacollectionendpoints.ParseDataCollectionEndpointIDInsensitively(*input); err == nil {
		nomalizedId := resourceId.ID()
		return &nomalizedId
	}

	return input
}
