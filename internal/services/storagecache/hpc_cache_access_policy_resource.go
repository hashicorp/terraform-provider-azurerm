// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/caches"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHPCCacheAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHPCCacheAccessPolicyCreateUpdate,
		Read:   resourceHPCCacheAccessPolicyRead,
		Update: resourceHPCCacheAccessPolicyCreateUpdate,
		Delete: resourceHPCCacheAccessPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CacheAccessPolicyID(id)
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
				ValidateFunc: validation.StringNotInSlice([]string{"default"}, false),
			},

			"hpc_cache_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: caches.ValidateCacheID,
			},

			"access_rule": {
				// Order doesn't matter for the access policies, as each one will be selected by one namespace path.
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 3,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"scope": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(caches.NfsAccessRuleScopeDefault),
								string(caches.NfsAccessRuleScopeNetwork),
								string(caches.NfsAccessRuleScopeHost),
							}, false),
						},

						"access": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(caches.NfsAccessRuleAccessRw),
								string(caches.NfsAccessRuleAccessRo),
								string(caches.NfsAccessRuleAccessNo),
							}, false),
						},

						"filter": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"suid_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"submount_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"root_squash_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"anonymous_uid": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"anonymous_gid": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},
		},
	}
}

func resourceHPCCacheAccessPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.Caches
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	cacheId, err := caches.ParseCacheID(d.Get("hpc_cache_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewCacheAccessPolicyID(cacheId.SubscriptionId, cacheId.ResourceGroupName, cacheId.CacheName, name)

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existCache, err := client.Get(ctx, *cacheId)
	if err != nil {
		return fmt.Errorf("checking for containing HPC Cache %s: %+v", cacheId, err)
	}

	m := existCache.Model
	if m == nil {
		return fmt.Errorf("nil Model of HPC Cache %s", cacheId)
	}

	prop := m.Properties
	if prop == nil {
		return fmt.Errorf("nil CacheProperties of HPC Cache %s", cacheId)
	}

	setting := prop.SecuritySettings
	if setting == nil {
		return fmt.Errorf("nil SecuritySettings of HPC Cache %s", cacheId)
	}

	policies := setting.AccessPolicies
	if policies == nil {
		return fmt.Errorf("nil AccessPolicies of HPC Cache %s", cacheId)
	}

	if d.IsNewResource() {
		p := CacheGetAccessPolicyByName(*policies, id.Name)
		if p != nil {
			return tf.ImportAsExistsError("azurerm_hpc_cache_access_policy", id.ID())
		}
	}

	p := caches.NfsAccessPolicy{
		Name:        id.Name,
		AccessRules: expandStorageCacheNfsAccessRules(d.Get("access_rule").(*pluginsdk.Set).List()),
	}

	*policies, err = CacheInsertOrUpdateAccessPolicy(*policies, p)
	if err != nil {
		return err
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *cacheId, *existCache.Model); err != nil {
		return fmt.Errorf("updating the HPC Cache for creating/updating Access Policy %q: %v", id, err)
	}

	d.SetId(id.ID())
	return resourceHPCCacheAccessPolicyRead(d, meta)
}

func resourceHPCCacheAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.Caches
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CacheAccessPolicyID(d.Id())
	if err != nil {
		return err
	}
	cacheId := caches.NewCacheID(id.SubscriptionId, id.ResourceGroup, id.CacheName)

	clearId := func(msg string) error {
		log.Printf("[DEBUG] %s - removing from state!", msg)
		d.SetId("")
		return nil
	}
	resp, err := client.Get(ctx, cacheId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return clearId(fmt.Sprintf("The containing HPC Cache %q was not found", cacheId))
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	m := resp.Model
	if m == nil {
		return fmt.Errorf("nil Model of HPC Cache %s", cacheId)
	}

	prop := m.Properties
	if prop == nil {
		return fmt.Errorf("nil CacheProperties of HPC Cache %s", cacheId)
	}

	setting := prop.SecuritySettings
	if setting == nil {
		return clearId(fmt.Sprintf("The containing HPC Cache %q has nil SecuritySettings", cacheId))
	}

	policies := setting.AccessPolicies
	if policies == nil {
		return clearId(fmt.Sprintf("The containing HPC Cache %q has nil AccessPolicies", cacheId))
	}

	p := CacheGetAccessPolicyByName(*policies, id.Name)
	if p == nil {
		return clearId(fmt.Sprintf("The %q was not found", id))
	}

	d.Set("name", id.Name)
	d.Set("hpc_cache_id", cacheId.ID())
	rules, err := flattenStorageCacheNfsAccessRules(p.AccessRules)
	if err != nil {
		return err
	}
	if err := d.Set("access_rule", rules); err != nil {
		return fmt.Errorf("setting `access_rule`: %v", err)
	}

	return nil
}

func resourceHPCCacheAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.Caches
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CacheAccessPolicyID(d.Id())
	if err != nil {
		return err
	}
	cacheId := caches.NewCacheID(id.SubscriptionId, id.ResourceGroup, id.CacheName)

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existCache, err := client.Get(ctx, cacheId)
	if err != nil {
		return fmt.Errorf("checking for containing HPC Cache %s: %+v", cacheId, err)
	}

	m := existCache.Model
	if m == nil {
		return nil
	}

	prop := m.Properties
	if prop == nil {
		return nil
	}

	settings := prop.SecuritySettings
	if settings == nil {
		return nil
	}

	policies := settings.AccessPolicies
	if policies == nil {
		return nil
	}

	*policies = CacheDeleteAccessPolicyByName(*policies, id.Name)

	if err = client.CreateOrUpdateThenPoll(ctx, cacheId, *existCache.Model); err != nil {
		return fmt.Errorf("updating the HPC Cache for deleting Access Policy %q: %v", id, err)
	}

	return nil
}
