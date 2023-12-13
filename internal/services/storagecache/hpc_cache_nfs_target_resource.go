// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/storagetargets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHPCCacheNFSTarget() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHPCCacheNFSTargetCreateOrUpdate,
		Update: resourceHPCCacheNFSTargetCreateOrUpdate,
		Read:   resourceHPCCacheNFSTargetRead,
		Delete: resourceHPCCacheNFSTargetDelete,

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

			"namespace_junction": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				// Confirmed with service team that they have a mac of 10 that is enforced by the backend.
				MaxItems: 10,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"namespace_path": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.CacheNamespacePath,
						},
						"nfs_export": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.CacheNFSExport,
						},
						"target_path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "",
							ValidateFunc: validate.CacheNFSTargetPath,
						},

						"access_policy_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "default",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"target_host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

			"verification_timer_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 31536000),
			},

			"write_back_timer_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 31536000),
			},
		},
	}
}

func resourceHPCCacheNFSTargetCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.StorageTargets
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure HPC Cache NFS Target creation.")
	id := storagetargets.NewStorageTargetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cache_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hpc_cache_nfs_target", id.ID())
		}
	}

	// Construct parameters
	param := storagetargets.StorageTarget{
		Properties: &storagetargets.StorageTargetProperties{
			Junctions:  expandNamespaceJunctions(d.Get("namespace_junction").(*pluginsdk.Set).List()),
			TargetType: storagetargets.StorageTargetTypeNfsThree,
			Nfs3: &storagetargets.Nfs3Target{
				Target:     pointer.To(d.Get("target_host_name").(string)),
				UsageModel: pointer.To(d.Get("usage_model").(string)),
			},
		},
	}

	if v, ok := d.GetOk("verification_timer_in_seconds"); ok {
		param.Properties.Nfs3.VerificationTimer = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("write_back_timer_in_seconds"); ok {
		param.Properties.Nfs3.WriteBackTimer = utils.Int64(int64(v.(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHPCCacheNFSTargetRead(d, meta)
}

func resourceHPCCacheNFSTargetRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] HPC Cache NFS Target %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving HPC Cache NFS Target %q: %+v", id, err)
	}

	d.Set("name", id.StorageTargetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cache_name", id.CacheName)

	if m := resp.Model; m != nil {
		if props := m.Properties; props != nil {
			if props.TargetType != storagetargets.StorageTargetTypeNfsThree {
				return fmt.Errorf("the type of this HPC Cache Target (%q) is not a NFS Target", id)
			}
			if nfs3 := props.Nfs3; nfs3 != nil {
				d.Set("target_host_name", nfs3.Target)
				d.Set("usage_model", nfs3.UsageModel)
				d.Set("verification_timer_in_seconds", pointer.From(nfs3.VerificationTimer))
				d.Set("write_back_timer_in_seconds", pointer.From(nfs3.WriteBackTimer))
			}
			if err := d.Set("namespace_junction", flattenNamespaceJunctions(props.Junctions)); err != nil {
				return fmt.Errorf(`error setting "namespace_junction"(%q): %+v`, id, err)
			}
		}
	}

	return nil
}

func resourceHPCCacheNFSTargetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.StorageTargets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagetargets.ParseStorageTargetID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id, storagetargets.DeleteOperationOptions{}); err != nil {
		return fmt.Errorf("deleting HPC Cache NFS Target (%q): %+v", id, err)
	}

	return nil
}

func expandNamespaceJunctions(input []interface{}) *[]storagetargets.NamespaceJunction {
	result := make([]storagetargets.NamespaceJunction, 0)

	for _, v := range input {
		b := v.(map[string]interface{})
		result = append(result, storagetargets.NamespaceJunction{
			NamespacePath:   pointer.To(b["namespace_path"].(string)),
			NfsExport:       pointer.To(b["nfs_export"].(string)),
			TargetPath:      pointer.To(b["target_path"].(string)),
			NfsAccessPolicy: pointer.To(b["access_policy_name"].(string)),
		})
	}

	return &result
}

func flattenNamespaceJunctions(input *[]storagetargets.NamespaceJunction) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, map[string]interface{}{
			"namespace_path":     pointer.From(e.NamespacePath),
			"nfs_export":         pointer.From(e.NfsExport),
			"target_path":        pointer.From(e.TargetPath),
			"access_policy_name": pointer.From(e.NfsAccessPolicy),
		})
	}

	return output
}
