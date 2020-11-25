package hpccache

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2020-03-01/storagecache"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHPCCacheNFSTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHPCCacheNFSTargetCreateOrUpdate,
		Update: resourceArmHPCCacheNFSTargetCreateOrUpdate,
		Read:   resourceArmHPCCacheNFSTargetRead,
		Delete: resourceArmHPCCacheNFSTargetDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parsers.StorageTargetID(id)
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
				ValidateFunc: validate.HPCCacheTargetName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cache_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"namespace_junction": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				// Confirmed with service team that they have a mac of 10 that is enforced by the backend.
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace_path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.HPCCacheNamespacePath,
						},
						"nfs_export": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.HPCCacheNFSExport,
						},
						"target_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "",
							ValidateFunc: validate.HPCCacheNFSTargetPath,
						},
					},
				},
			},

			"target_host_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"usage_model": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"READ_HEAVY_INFREQ",
					"WRITE_WORKLOAD_15",
					"WRITE_AROUND",
				}, false),
			},
		},
	}
}

func resourceArmHPCCacheNFSTargetCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure HPC Cache NFS Target creation.")
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	cache := d.Get("cache_name").(string)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, cache, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for existing HPC Cache NFS Target %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_hpc_cache_nfs_target", *resp.ID)
		}
	}

	// Construct parameters
	param := &storagecache.StorageTarget{
		BasicStorageTargetProperties: &storagecache.Nfs3TargetProperties{
			Junctions:  expandNamespaceJunctions(d.Get("namespace_junction").(*schema.Set).List()),
			TargetType: storagecache.TargetTypeNfs3,
			Nfs3: &storagecache.Nfs3Target{
				Target:     utils.String(d.Get("target_host_name").(string)),
				UsageModel: utils.String(d.Get("usage_model").(string)),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, cache, name, param)
	if err != nil {
		return fmt.Errorf("Error creating HPC Cache NFS Target %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of HPC Cache NFS Target %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, cache, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HPC Cache NFS Target %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Error retrieving HPC Cache NFS Target %q (Resource Group %q): `id` was nil", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHPCCacheNFSTargetRead(d, meta)
}

func resourceArmHPCCacheNFSTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.StorageTargetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Cache, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HPC Cache NFS Target %q was not found (Resource Group %q, Cahe %q) - removing from state!", id.Name, id.ResourceGroup, id.Cache)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving HPC Cache NFS Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.Cache, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cache_name", id.Cache)

	if props := resp.BasicStorageTargetProperties; props != nil {
		props, ok := props.AsNfs3TargetProperties()
		if !ok {
			return fmt.Errorf("The type of this HPC Cache Target %q (Resource Group %q, Cahe %q) is not a NFS Target", id.Name, id.ResourceGroup, id.Cache)
		}
		if nfs3 := props.Nfs3; nfs3 != nil {
			d.Set("target_host_name", nfs3.Target)
			d.Set("usage_model", nfs3.UsageModel)
		}
		if err := d.Set("namespace_junction", flattenNamespaceJunctions(props.Junctions)); err != nil {
			return fmt.Errorf(`Error setting "namespace_junction" %q (Resource Group %q, Cahe %q): %+v`, id.Name, id.ResourceGroup, id.Cache, err)
		}
	}

	return nil
}

func resourceArmHPCCacheNFSTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.StorageTargetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.StorageTargetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Cache, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting HPC Cache NFS Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.Cache, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of HPC Cache NFS Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.Cache, err)
	}

	return nil
}

func expandNamespaceJunctions(input []interface{}) *[]storagecache.NamespaceJunction {
	result := make([]storagecache.NamespaceJunction, 0)

	for _, v := range input {
		b := v.(map[string]interface{})
		result = append(result, storagecache.NamespaceJunction{
			NamespacePath: utils.String(b["namespace_path"].(string)),
			NfsExport:     utils.String(b["nfs_export"].(string)),
			TargetPath:    utils.String(b["target_path"].(string)),
		})
	}

	return &result
}

func flattenNamespaceJunctions(input *[]storagecache.NamespaceJunction) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		namespacePath := ""
		if v := e.NamespacePath; v != nil {
			namespacePath = *v
		}

		nfsExport := ""
		if v := e.NfsExport; v != nil {
			nfsExport = *v
		}

		targetPath := ""
		if v := e.TargetPath; v != nil {
			targetPath = *v
		}

		output = append(output, map[string]interface{}{
			"namespace_path": namespacePath,
			"nfs_export":     nfsExport,
			"target_path":    targetPath,
		})
	}

	return output
}
