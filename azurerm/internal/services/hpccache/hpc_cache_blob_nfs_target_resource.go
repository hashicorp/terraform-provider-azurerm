package hpccache

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-03-01/storagecache"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/validate"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceHPCCacheBlobNFSTarget() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHPCCacheBlobNFSTargetCreateUpdate,
		Read:   resourceHPCCacheBlobNFSTargetRead,
		Update: resourceHPCCacheBlobNFSTargetCreateUpdate,
		Delete: resourceHPCCacheBlobNFSTargetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StorageTargetID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				ValidateFunc: storageValidate.StorageContainerResourceManagerID,
			},

			// TODO: use SDK enums once following issue is addressed
			// https://github.com/Azure/azure-rest-api-specs/issues/13839
			"usage_model": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"READ_HEAVY_INFREQ",
					"READ_HEAVY_CHECK_180",
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
	id := parse.NewStorageTargetID(subscriptionId, resourceGroup, cache, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.CacheName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_hpc_cache_blob_nfs_target", id.ID())
		}
	}

	namespacePath := d.Get("namespace_path").(string)
	containerId := d.Get("storage_container_id").(string)

	// Construct parameters
	namespaceJunction := []storagecache.NamespaceJunction{
		{
			NamespacePath:   &namespacePath,
			TargetPath:      utils.String("/"),
			NfsExport:       utils.String("/"),
			NfsAccessPolicy: utils.String(d.Get("access_policy_name").(string)),
		},
	}
	param := &storagecache.StorageTarget{
		StorageTargetProperties: &storagecache.StorageTargetProperties{
			Junctions:  &namespaceJunction,
			TargetType: storagecache.StorageTargetTypeBlobNfs,
			BlobNfs: &storagecache.BlobNfsTarget{
				Target:     utils.String(containerId),
				UsageModel: utils.String(d.Get("usage_model").(string)),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.CacheName, id.Name, param)
	if err != nil {
		return fmt.Errorf("Error creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHPCCacheBlobNFSTargetRead(d, meta)
}

func resourceHPCCacheBlobNFSTargetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageTargetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.CacheName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cache_name", id.CacheName)

	if props := resp.StorageTargetProperties; props != nil {
		if props.TargetType != storagecache.StorageTargetTypeBlobNfs {
			return fmt.Errorf("The type of this HPC Cache Target %s is not a Blob NFS Target", id)
		}

		storageContainerId := ""
		usageModel := ""
		if b := props.BlobNfs; b != nil {
			if b.Target != nil {
				storageContainerId = *b.Target
			}
			if b.UsageModel != nil {
				usageModel = *b.UsageModel
			}
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
	return nil
}

func resourceHPCCacheBlobNFSTargetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageTargetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.CacheName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
