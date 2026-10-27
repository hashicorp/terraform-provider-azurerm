// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/resources/redis_cache/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-name basicWithSSL

var redisCacheResourceName = "azurerm_redis_cache"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceRedisCacheCreate,
		Read:     resourceRedisCacheRead,
		Update:   resourceRedisCacheUpdate,
		Delete:   resourceRedisCacheDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&redisresources.RediId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&redisresources.RediId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(180 * time.Minute),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RedisCacheV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: redisCacheSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("sku_name", func(ctx context.Context, old, new, meta any) bool {
				// downgrade the SKU is not supported, recreate the resource
				if old.(string) != "" && new.(string) != "" {
					return skuWeight[old.(string)] > skuWeight[new.(string)]
				}
				return false
			}),
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				// Entra (AD) auth has to be set to disable access keys auth
				// https://learn.microsoft.com/azure/azure-cache-for-redis/cache-azure-active-directory-for-authentication

				accessKeysAuthenticationEnabled := diff.Get("access_keys_authentication_enabled").(bool)
				activeDirectoryAuthenticationEnabled := diff.Get("redis_configuration.0.active_directory_authentication_enabled").(bool)

				log.Printf("[DEBUG] CustomizeDiff: access_keys_authentication_enabled: %v, active_directory_authentication_enabled: %v", accessKeysAuthenticationEnabled, activeDirectoryAuthenticationEnabled)

				if !accessKeysAuthenticationEnabled && !activeDirectoryAuthenticationEnabled {
					return fmt.Errorf("`active_directory_authentication_enabled` must be enabled in order to disable `access_keys_authentication_enabled`")
				}

				return nil
			}),
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				// Replicates validation rules from Azure CLI
				// https://github.com/Azure/azure-cli/blob/131634d374fe704920862a2e5b0745e61af9bc89/src/azure-cli/azure/cli/command_modules/redis/custom.py#L13
				skuName := diff.Get("sku_name").(string)
				family := diff.Get("family").(string)
				capacity := diff.Get("capacity").(int)
				validCombinations := map[string][]string{
					string(redisresources.SkuNameBasic):    {"C0", "C1", "C2", "C3", "C4", "C5", "C6"},
					string(redisresources.SkuNameStandard): {"C0", "C1", "C2", "C3", "C4", "C5", "C6"},
					string(redisresources.SkuNamePremium):  {"P1", "P2", "P3", "P4", "P5"},
				}

				familyCapacity := fmt.Sprintf("%s%d", strings.ToUpper(family), capacity)
				if !slices.Contains(validCombinations[skuName], familyCapacity) {
					return fmt.Errorf("invalid combination of `sku_name`, `family`, and `capacity`: '%[1]s: %[2]s'. Valid combinations for '%[1]s' are: %[3]v", skuName, familyCapacity, validCombinations[skuName])
				}

				return nil
			}),
		),
	}
}
