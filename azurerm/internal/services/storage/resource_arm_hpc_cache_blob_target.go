package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2019-11-01/storagecache"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	storage "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
)

func resourceArmHPCCacheBlobTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHPCCacheBlobTargetCreateOrUpdate,
		Update: resourceArmHPCCacheBlobTargetCreateOrUpdate,
		Read:   resourceArmHPCCacheBlobTargetRead,
		Delete: resourceArmHPCCacheBlobTargetDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parsers.HPCCacheTargetID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.HPCCacheTargetName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cache_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"namespace_path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: storageValidate.HPCCacheNamespacePath,
			},

			"storage_container_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmHPCCacheBlobTargetCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StorageTargetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure HPC Cache Blob Target creation.")
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	cache := d.Get("cache_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	// We need to convert storage container id from the form used in side hashicorp azure storage sdk
	// to the form used in original azure storage sdk, which is expected by storage target client.

	azureContainerId, err := toAzureStorageContianerID(
		containerId,
		meta.(*clients.Client).Account.SubscriptionId,
		resourceGroup,
	)
	if err != nil {
		return fmt.Errorf("Error converting container id to azure container id (HPC Cache Blob Target %q, Resource Group %q, Cahe %q): %w", name, resourceGroup, cache, err)
	}

	// Construct parameters
	namespaceJunction := []storagecache.NamespaceJunction{
		{
			NamespacePath: &namespacePath,
			TargetPath:    utils.String("/"),
		},
	}
	param := &storagecache.StorageTarget{
		StorageTargetProperties: &storagecache.StorageTargetProperties{
			Junctions:  &namespaceJunction,
			TargetType: storagecache.StorageTargetTypeClfs,
			Clfs: &storagecache.ClfsTarget{
				Target: &azureContainerId,
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

	return resourceArmHPCCacheBlobTargetRead(d, meta)
}

func resourceArmHPCCacheBlobTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StorageTargetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.HPCCacheTargetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Cache, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HPC Cache Blob Target %q was not found (Resource Group %q, Cahe %q) - removing from state!", id.Name, id.ResourceGroup, id.Cache)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving HPC Cache Blob Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.Cache, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cache_name", id.Cache)

	if resp.StorageTargetProperties == nil {
		return fmt.Errorf("Error retrieving HPC Cache Blob Target %q (Resource Group %q, Cahe %q): `properties` was nil", id.Name, id.ResourceGroup, id.Cache)
	}
	props := *resp.StorageTargetProperties

	containerId := ""
	if props.Clfs != nil && props.Clfs.Target != nil {
		// Convert container id from azure form to the form used in current
		// storage sdk package.
		azureContainerId := *props.Clfs.Target
		var err error
		containerId, err = fromAzureStorageContainerID(ctx, meta.(*clients.Client).Storage, azureContainerId)
		if err != nil {
			return fmt.Errorf("Error converting container id from azure container id (HPC Cache Blob Target %q, Resource Group %q, Cahe %q): %w", id.Name, id.ResourceGroup, id.Cache, err)
		}
	}
	d.Set("storage_container_id", containerId)

	namespacePath := ""
	// There is only one namespace path allowed for blob container storage target,
	// which maps to the root path of it.
	if props.Junctions != nil && len(*props.Junctions) == 1 && (*props.Junctions)[0].NamespacePath != nil {
		namespacePath = *(*props.Junctions)[0].NamespacePath
	}
	d.Set("namespace_path", namespacePath)

	return nil
}

func resourceArmHPCCacheBlobTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StorageTargetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.HPCCacheTargetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Cache, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting HPC Cache Blob Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.Cache, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error wating for deletiion of HPC Cache Blob Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.Cache, err)
	}

	return nil
}

// toAzureStorageContianerID convert container id from form used in "giovanni"
// package to form used in azure.
func toAzureStorageContianerID(id, subid, resgroup string) (string, error) {
	idinfo, err := containers.ParseResourceID(id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s",
		subid, resgroup, idinfo.AccountName, idinfo.ContainerName), nil
}

// fromAzureStorageContainerID convert container id from form used in azure
// to form used in "giovanni" package.
func fromAzureStorageContainerID(ctx context.Context, storageClient *storage.Client, input string) (string, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return "", err
	}
	var accountName, containerName string

	if accountName, err = id.PopSegment("storageAccounts"); err != nil {
		return "", err
	}

	if containerName, err = id.PopSegment("containers"); err != nil {
		return "", err
	}

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return "", fmt.Errorf("Error retrieving Account %q for Container %q: %s", accountName, containerName, err)
	}
	if account == nil {
		return "", fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return "", fmt.Errorf("Error building Containers Client: %s", err)
	}

	return client.GetResourceID(accountName, containerName), nil
}
