package hpccache

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-03-01/storagecache"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceHPCCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceHPCCacheCreateOrUpdate,
		Update: resourceHPCCacheCreateOrUpdate,
		Read:   resourceHPCCacheRead,
		Delete: resourceHPCCacheDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CacheID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cache_size_in_gb": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					3072,
					6144,
					12288,
					24576,
					49152,
				}),
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard_2G",
					"Standard_4G",
					"Standard_8G",
				}, false),
			},

			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1500,
				ValidateFunc: validation.IntBetween(576, 1500),
			},

			// TODO 3.0: remove this property
			"root_squash_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				// TODO 3.0: remove "Computed: true" and add "Default: true"
				// The old resource has no consistent default for the rootSquash setting. In order not to
				// break users, we intentionally mark this property as Computed.
				// https://docs.microsoft.com/en-us/azure/hpc-cache/configuration#configure-root-squash.
				Computed:   true,
				Deprecated: "This is deprecated in favor of `default_access_policy.0.access_rule.x.root_squash_enabled`, where the scope of access_rule is `default`. Will be removed in v3.0",
			},

			"default_access_policy": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				// This is computed because there is always a "default" policy in the cache. It is created together with the cache, and users can't remove it.
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_rule": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(storagecache.Default),
											string(storagecache.Network),
											string(storagecache.Host),
										}, false),
									},

									"access": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(storagecache.NfsAccessRuleAccessRw),
											string(storagecache.NfsAccessRuleAccessRo),
											string(storagecache.NfsAccessRuleAccessNo),
										}, false),
									},

									"filter": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"suid_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"submount_access_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"root_squash_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"anonymous_uid": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},

									"anonymous_gid": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
								},
							},
						},
					},
				},
			},

			"mount_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceHPCCacheCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.CachesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure HPC Cache creation.")
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewCacheID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing HPC Cache %q: %s", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_hpc_cache", id.ID())
		}
	}

	location := d.Get("location").(string)
	cacheSize := d.Get("cache_size_in_gb").(int)
	subnet := d.Get("subnet_id").(string)
	skuName := d.Get("sku_name").(string)
	mtu := d.Get("mtu").(int)

	var accessPolicies []storagecache.NfsAccessPolicy
	if !d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving existing HPC Cache %q: %v", id, err)
		}
		if prop := existing.CacheProperties; prop != nil {
			if settings := existing.SecuritySettings; settings != nil {
				if policies := settings.AccessPolicies; policies != nil {
					accessPolicies = *policies
				}
			}
		}
	}
	defaultAccessPolicy, err := expandStorageCacheDefaultAccessPolicy(d.Get("default_access_policy").([]interface{}), d.Get("root_squash_enabled").(bool))
	if err != nil {
		return err
	}
	accessPolicies, err = CacheInsertOrUpdateAccessPolicy(accessPolicies, defaultAccessPolicy)
	if err != nil {
		return err
	}

	cache := &storagecache.Cache{
		Name:     utils.String(name),
		Location: utils.String(location),
		CacheProperties: &storagecache.CacheProperties{
			CacheSizeGB: utils.Int32(int32(cacheSize)),
			Subnet:      utils.String(subnet),
			NetworkSettings: &storagecache.CacheNetworkSettings{
				Mtu: utils.Int32(int32(mtu)),
			},
			SecuritySettings: &storagecache.CacheSecuritySettings{
				AccessPolicies: &accessPolicies,
			},
		},
		Sku: &storagecache.CacheSku{
			Name: utils.String(skuName),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, cache)
	if err != nil {
		return fmt.Errorf("Error creating HPC Cache %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for HPC Cache %q (Resource Group %q) to finish provisioning: %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HPC Cache %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for HPC Cache %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceHPCCacheRead(d, meta)
}

func resourceHPCCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.CachesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CacheID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HPC Cache %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HPC Cache %q: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", resp.Location)

	if props := resp.CacheProperties; props != nil {
		d.Set("cache_size_in_gb", props.CacheSizeGB)
		d.Set("subnet_id", props.Subnet)
		d.Set("mount_addresses", utils.FlattenStringSlice(props.MountAddresses))
		if props.NetworkSettings != nil {
			d.Set("mtu", props.NetworkSettings.Mtu)
		}
		if securitySettings := props.SecuritySettings; securitySettings != nil {
			if securitySettings.AccessPolicies != nil {
				defaultPolicy := CacheGetAccessPolicyByName(*securitySettings.AccessPolicies, "default")
				if defaultPolicy != nil {
					defaultAccessPolicy, err := flattenStorageCacheNfsDefaultAccessPolicy(*defaultPolicy)
					if err != nil {
						return err
					}
					if err := d.Set("default_access_policy", defaultAccessPolicy); err != nil {
						return fmt.Errorf("setting `default_access_policy`: %v", err)
					}

					// Set the "root_squash_enabled" for backward compatibility.
					// TODO 3.0 - remove this part.
					deprecatedRootSquashEnabled := false
					if defaultPolicy.AccessRules != nil {
						accessRule, ok := CacheGetAccessPolicyRuleByScope(*defaultPolicy.AccessRules, storagecache.Default)
						if ok && accessRule.RootSquash != nil {
							deprecatedRootSquashEnabled = *accessRule.RootSquash
						}
					}
					d.Set("root_squash_enabled", deprecatedRootSquashEnabled)
				}
			}
		}
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	return nil
}

func resourceHPCCacheDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HPCCache.CachesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CacheID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting HPC Cache %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of HPC Cache %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandStorageCacheDefaultAccessPolicy(input []interface{}, deprecatedRootSquashed bool) (storagecache.NfsAccessPolicy, error) {
	// Use the deprecated "root_squashed_enabled" property to setup the default scoped access policy for backward compatibility.
	if len(input) == 0 {
		accessRules := []storagecache.NfsAccessRule{
			{
				Scope:          storagecache.Default,
				Access:         storagecache.NfsAccessRuleAccessRw,
				Suid:           utils.Bool(false),
				SubmountAccess: utils.Bool(true),
				RootSquash:     utils.Bool(false),
			},
		}

		return storagecache.NfsAccessPolicy{
			Name:        utils.String("default"),
			AccessRules: &accessRules,
		}, nil
	}

	if deprecatedRootSquashed {
		return storagecache.NfsAccessPolicy{}, fmt.Errorf("`root_squashed_enabled` can't be used together with `default_access_policy`, prefer using the latter one exclusively")
	}

	return storagecache.NfsAccessPolicy{
		Name:        utils.String("default"),
		AccessRules: expandStorageCacheNfsAccessRules(input[0].(map[string]interface{})["access_rule"].(*schema.Set).List()),
	}, nil
}

func flattenStorageCacheNfsDefaultAccessPolicy(input storagecache.NfsAccessPolicy) ([]interface{}, error) {
	rules, err := flattenStorageCacheNfsAccessRules(input.AccessRules)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		map[string]interface{}{
			"access_rule": rules,
		},
	}, nil
}

func expandStorageCacheNfsAccessRules(input []interface{}) *[]storagecache.NfsAccessRule {
	var out []storagecache.NfsAccessRule
	for _, accessRuleRaw := range input {
		b := accessRuleRaw.(map[string]interface{})
		out = append(out, storagecache.NfsAccessRule{
			Scope:          storagecache.NfsAccessRuleScope(b["scope"].(string)),
			Access:         storagecache.NfsAccessRuleAccess(b["access"].(string)),
			Filter:         utils.String(b["filter"].(string)),
			Suid:           utils.Bool(b["suid_enabled"].(bool)),
			SubmountAccess: utils.Bool(b["submount_access_enabled"].(bool)),
			RootSquash:     utils.Bool(b["root_squash_enabled"].(bool)),
			AnonymousUID:   utils.String(strconv.Itoa(b["anonymous_uid"].(int))),
			AnonymousGID:   utils.String(strconv.Itoa(b["anonymous_gid"].(int))),
		})
	}
	return &out
}

func flattenStorageCacheNfsAccessRules(input *[]storagecache.NfsAccessRule) ([]interface{}, error) {
	if input == nil {
		return nil, nil
	}

	var rules []interface{}
	for _, accessRule := range *input {
		filter := ""
		if accessRule.Filter != nil {
			filter = *accessRule.Filter
		}

		suidEnabled := false
		if accessRule.Suid != nil {
			suidEnabled = *accessRule.Suid
		}

		submountAccessEnabled := false
		if accessRule.SubmountAccess != nil {
			submountAccessEnabled = *accessRule.SubmountAccess
		}

		rootSquashEnabled := false
		if accessRule.RootSquash != nil {
			rootSquashEnabled = *accessRule.RootSquash
		}

		anonymousUID := 0
		if accessRule.AnonymousUID != nil {
			var err error
			anonymousUID, err = strconv.Atoi(*accessRule.AnonymousUID)
			if err != nil {
				return nil, fmt.Errorf("converting `anonymous_uid` from string to int")
			}
		}

		anonymousGID := 0
		if accessRule.AnonymousGID != nil {
			var err error
			anonymousGID, err = strconv.Atoi(*accessRule.AnonymousGID)
			if err != nil {
				return nil, fmt.Errorf("converting `anonymous_gid` from string to int")
			}
		}

		rules = append(rules, map[string]interface{}{
			"scope":                   accessRule.Scope,
			"access":                  accessRule.Access,
			"filter":                  filter,
			"suid_enabled":            suidEnabled,
			"submount_access_enabled": submountAccessEnabled,
			"root_squash_enabled":     rootSquashEnabled,
			"anonymous_uid":           anonymousUID,
			"anonymous_gid":           anonymousGID,
		})
	}

	return rules, nil
}
