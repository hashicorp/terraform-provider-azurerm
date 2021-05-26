package hpccache

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-03-01/storagecache"
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

func resourceHPCCacheBlobTarget() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHPCCacheBlobTargetCreateOrUpdate,
		Update: resourceHPCCacheBlobTargetCreateOrUpdate,
		Read:   resourceHPCCacheBlobTargetRead,
		Delete: resourceHPCCacheBlobTargetDelete,

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

			"cache_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure HPC Cache Blob Target creation.")
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	cache := d.Get("cache_name").(string)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, cache, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for existing HPC Cache Blob Target %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_hpc_cache_blob_target", *resp.ID)
		}
	}

	namespacePath := d.Get("namespace_path").(string)
	containerId := d.Get("storage_container_id").(string)

	// Construct parameters
	namespaceJunction := []storagecache.NamespaceJunction{
		{
			NamespacePath:   &namespacePath,
			TargetPath:      utils.String("/"),
			NfsAccessPolicy: utils.String(d.Get("access_policy_name").(string)),
		},
	}
	param := &storagecache.StorageTarget{
		StorageTargetProperties: &storagecache.StorageTargetProperties{
			Junctions:  &namespaceJunction,
			TargetType: storagecache.StorageTargetTypeClfs,
			Clfs: &storagecache.ClfsTarget{
				Target: utils.String(containerId),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, cache, name, param)
	if err != nil {
		return fmt.Errorf("Error creating HPC Cache Blob Target %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of HPC Cache Blob Target %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, cache, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HPC Cache Blob Target %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Error retrieving HPC Cache Blob Target %q (Resource Group %q): `id` was nil", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceHPCCacheBlobTargetRead(d, meta)
}

func resourceHPCCacheBlobTargetRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] HPC Cache Blob Target %q was not found (Resource Group %q, Cache %q) - removing from state!", id.Name, id.ResourceGroup, id.CacheName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HPC Cache Blob Target %q (Resource Group %q, Cache %q): %+v", id.Name, id.ResourceGroup, id.CacheName, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cache_name", id.CacheName)

	if props := resp.StorageTargetProperties; props != nil {
		if props.TargetType != storagecache.StorageTargetTypeClfs {
			return fmt.Errorf("The type of this HPC Cache Target %q (Resource Group %q, Cahe %q) is not a Blob Target", id.Name, id.ResourceGroup, id.CacheName)
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

	return nil
}

func resourceHPCCacheBlobTargetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageTargetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.CacheName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting HPC Cache Blob Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.CacheName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of HPC Cache Blob Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.CacheName, err)
	}

	return nil
}
