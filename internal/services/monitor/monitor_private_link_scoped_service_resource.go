// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameOptional(),

			"scope_resource_group_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"resource_group_name", "scope_resource_group_name", "scope_resource_id"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"scope_subscription_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"scope_name", "scope_resource_id"},
				ValidateFunc: validate.PrivateLinkScopeName,
			},

			"scope_resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"scope_name", "scope_resource_id"},
				ValidateFunc: privatelinkscopesapis.ValidatePrivateLinkScopeID,
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
	resourceGroupName := d.Get("resource_group_name").(string)
	scopeName := d.Get("scope_name").(string)
	if scopeResourceGroupName, ok := d.GetOk("scope_resource_group_name"); ok {
		resourceGroupName = scopeResourceGroupName.(string)
	}
	if scopeSubscriptionId, ok := d.GetOk("scope_subscription_id"); ok {
		subscriptionId = scopeSubscriptionId.(string)
	}
	client := meta.(*clients.Client).Monitor.PrivateLinkScopedResourcesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if scopeIdRaw, ok := d.GetOk("scope_resource_id"); ok {
		scopeId, err := privatelinkscopesapis.ParsePrivateLinkScopeID(scopeIdRaw.(string))
		if err != nil {
			return fmt.Errorf("parsing scope ID: %+v", err)
		}

		subscriptionId = scopeId.SubscriptionId
		resourceGroupName = scopeId.ResourceGroupName
		scopeName = scopeId.PrivateLinkScopeName
	}

	id := privatelinkscopedresources.NewScopedResourceID(subscriptionId, resourceGroupName, scopeName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			LinkedResourceId: pointer.To(d.Get("linked_resource_id").(string)),
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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
	if scopeIdRaw, ok := d.GetOk("scope_resource_id"); ok {
		scopeId, err := privatelinkscopesapis.ParsePrivateLinkScopeID(scopeIdRaw.(string))
		if err != nil {
			return fmt.Errorf("parsing scope ID: %+v", err)
		}
		d.Set("scope_resource_id", scopeId.ID())
	} else {
		d.Set("scope_name", id.PrivateLinkScopeName)
		if _, ok := d.GetOk("scope_resource_group_name"); ok {
			d.Set("scope_resource_group_name", id.ResourceGroupName)
		} else {
			d.Set("resource_group_name", id.ResourceGroupName)
		}
		if _, ok := d.GetOk("scope_subscription_id"); ok {
			d.Set("scope_subscription_id", id.SubscriptionId)
		}
	}

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

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
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
