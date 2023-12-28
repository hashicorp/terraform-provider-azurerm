// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/storagetargets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHPCCacheBlobTarget() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHPCCacheBlobTargetCreateOrUpdate,
		Update: resourceHPCCacheBlobTargetCreateOrUpdate,
		Read:   resourceHPCCacheBlobTargetRead,
		Delete: resourceHPCCacheBlobTargetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := storagetargets.ParseStorageTargetID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageTargetName,
			},

			"cache_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"namespace_path": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CacheNamespacePath,
			},

			"storage_container_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageContainerID,
			},

			"access_policy_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "default",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceHPCCacheBlobTargetCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.StorageTargets
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure HPC Cache Blob Target creation.")
	id := storagetargets.NewStorageTargetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cache_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hpc_cache_blob_target", id.ID())
		}
	}

	namespacePath := d.Get("namespace_path").(string)
	containerId := d.Get("storage_container_id").(string)

	// Construct parameters
	param := storagetargets.StorageTarget{
		Properties: &storagetargets.StorageTargetProperties{
			Junctions: &[]storagetargets.NamespaceJunction{
				{
					NamespacePath:   &namespacePath,
					TargetPath:      pointer.To("/"),
					NfsAccessPolicy: pointer.To(d.Get("access_policy_name").(string)),
				},
			},
			TargetType: storagetargets.StorageTargetTypeClfs,
			Clfs: &storagetargets.ClfsTarget{
				Target: pointer.To(containerId),
			},
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHPCCacheBlobTargetRead(d, meta)
}

func resourceHPCCacheBlobTargetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.StorageTargets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagetargets.ParseStorageTargetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] HPC Cache Blob Target %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HPC Cache Blob Target %q: %+v", id, err)
	}

	d.Set("name", id.StorageTargetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cache_name", id.CacheName)

	if m := resp.Model; m != nil {
		if props := m.Properties; props != nil {
			if props.TargetType != storagetargets.StorageTargetTypeClfs {
				return fmt.Errorf("The type of this HPC Cache Target %q is not a Blob Target", id)
			}

			storageContainerId := ""
			if clfs := props.Clfs; clfs != nil && clfs.Target != nil {
				storageContainerId = *clfs.Target
			}
			d.Set("storage_container_id", storageContainerId)

			namespacePath := ""
			accessPolicy := ""
			// There is only one namespace path allowed for blob container storage target,
			// which maps to the root path of it.
			if props.Junctions != nil && len(*props.Junctions) == 1 && (*props.Junctions)[0].NamespacePath != nil {
				namespacePath = *(*props.Junctions)[0].NamespacePath
				accessPolicy = *(*props.Junctions)[0].NfsAccessPolicy
			}
			d.Set("namespace_path", namespacePath)
			d.Set("access_policy_name", accessPolicy)
		}
	}

	return nil
}

func resourceHPCCacheBlobTargetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.StorageTargets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagetargets.ParseStorageTargetID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id, storagetargets.DeleteOperationOptions{}); err != nil {
		return fmt.Errorf("deleting HPC Cache Blob Target %q : %+v", id, err)
	}

	return nil
}
