// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name resource_group -service-package-name resource -properties "name" -known-values "subscription_id:data.Subscriptions.Primary"

func resourceResourceGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceResourceGroupCreate,
		Read:     resourceResourceGroupRead,
		Update:   resourceResourceGroupUpdate,
		Delete:   resourceResourceGroupDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.ResourceGroupId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"tags": commonschema.Tags(),

			"managed_by": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.ResourceGroupId{}),
		},
	}
}

func resourceResourceGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewResourceGroupID(meta.(*clients.Client).Account.SubscriptionId, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_resource_group", id.ID())
	}

	parameters := resourcegroups.ResourceGroup{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v := d.Get("managed_by").(string); v != "" {
		parameters.ManagedBy = pointer.To(v)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	// custom poller to account for replication delays in the eventual consistency responses of newly created RG resources
	pollerType := custompollers.NewResourceGroupCreatePoller(client, id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err = poller.PollUntilDone(ctx); err != nil {
		return err
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceResourceGroupRead(d, meta)
}

func resourceResourceGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	patch := resourcegroups.ResourceGroupPatchable{}

	if d.HasChange("managed_by") {
		patch.ManagedBy = pointer.To(d.Get("managed_by").(string))
	}

	if d.HasChange("tags") {
		patch.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, patch); err != nil {
		return fmt.Errorf("updating %q: %+v", *id, err)
	}

	return resourceResourceGroupRead(d, meta)
}

func resourceResourceGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading resource group %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading resource group: %+v", err)
	}

	d.Set("name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("managed_by", pointer.From(model.ManagedBy))
		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceResourceGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	// conditionally check for nested resources and error if they exist
	if meta.(*clients.Client).Features.ResourceGroup.PreventDeletionIfContainsResources {
		// Resource groups sometimes hold on to resource information after the resources have been deleted. We'll retry this check to account for that eventual consistency.
		deletePollerContext, deletePollerCancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer deletePollerCancel()

		pollerType := custompollers.NewResourceGroupPreventDeletePoller(client, *id)
		poller := pollers.NewRetryOnErrorPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow, true)
		if err := poller.PollUntilDone(deletePollerContext); err != nil {
			return err
		}
	}

	// This is not the usual pattern for destroys on go-azure-sdk, however, this functionally the same as the resource
	// worked before refactoring, so behaviour has been maintained. This should be investigated in future and brought
	// in-line if possible.
	if resp, err := client.Delete(ctx, *id, resourcegroups.DefaultDeleteOperationOptions()); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	} else {
		if err := resp.Poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling deleting %s: %+v", *id, err)
		}
	}

	return nil
}
