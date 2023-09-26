// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hpccache

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-01-01/storagetargets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hpccache/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHPCCacheBlobNFSTarget() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHPCCacheBlobNFSTargetCreateUpdate,
		Read:   resourceHPCCacheBlobNFSTargetRead,
		Update: resourceHPCCacheBlobNFSTargetCreateUpdate,
		Delete: resourceHPCCacheBlobNFSTargetDelete,

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

			"resource_group_name": commonschema.ResourceGroupName(),

			"cache_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

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

			// TODO: use SDK enums once following issue is addressed
			// https://github.com/Azure/azure-rest-api-specs/issues/13839
			"usage_model": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"READ_HEAVY_INFREQ",
					"READ_HEAVY_CHECK_180",
					"READ_ONLY",
					"READ_WRITE",
					"WRITE_WORKLOAD_15",
					"WRITE_AROUND",
					"WRITE_WORKLOAD_CHECK_30",
					"WRITE_WORKLOAD_CHECK_60",
					"WRITE_WORKLOAD_CLOUDWS",
				}, false),
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

func resourceHPCCacheBlobNFSTargetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	cache := d.Get("cache_name").(string)
	id := storagetargets.NewStorageTargetID(subscriptionId, resourceGroup, cache, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hpc_cache_blob_nfs_target", id.ID())
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
					NfsExport:       pointer.To("/"),
					NfsAccessPolicy: pointer.To(d.Get("access_policy_name").(string)),
				},
			},
			TargetType: storagetargets.StorageTargetTypeBlobNfs,
			BlobNfs: &storagetargets.BlobNfsTarget{
				Target:     pointer.To(containerId),
				UsageModel: pointer.To(d.Get("usage_model").(string)),
			},
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHPCCacheBlobNFSTargetRead(d, meta)
}

func resourceHPCCacheBlobNFSTargetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagetargets.ParseStorageTargetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.StorageTargetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cache_name", id.CacheName)

	if m := resp.Model; m != nil {
		if props := m.Properties; props != nil {
			if props.TargetType != storagetargets.StorageTargetTypeBlobNfs {
				return fmt.Errorf("The type of this HPC Cache Target %s is not a Blob NFS Target", id)
			}

			storageContainerId := ""
			usageModel := ""
			if b := props.BlobNfs; b != nil {
				storageContainerId = pointer.From(b.Target)
				usageModel = pointer.From(b.UsageModel)
			}
			d.Set("storage_container_id", storageContainerId)
			d.Set("usage_model", usageModel)

			namespacePath := ""
			accessPolicy := ""
			// There is only one namespace path allowed for the blob nfs target,
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

func resourceHPCCacheBlobNFSTargetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagetargets.ParseStorageTargetID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id, storagetargets.DeleteOperationOptions{}); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
