// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryCreate,
		Read:   resourceContainerRegistryRead,
		Update: resourceContainerRegistryUpdate,
		Delete: resourceContainerRegistryDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RegistryV0ToV1{},
			1: migration.RegistryV1ToV2{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := registries.ParseRegistryID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: containerRegistrySchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			sku := d.Get("sku").(string)

			geoReplications := d.Get("georeplications").([]any)
			// if locations have been specified for geo-replication then, the SKU has to be Premium
			if len(geoReplications) > 0 && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return errors.New("an ACR geo-replication can only be applied when using the Premium Sku")
			}

			// ensure location is different than any location of the geo-replication
			var geoReplicationLocations []string
			for _, v := range geoReplications {
				v := v.(map[string]any)
				geoReplicationLocations = append(geoReplicationLocations, location.Normalize(v["location"].(string)))
			}
			location := location.Normalize(d.Get("location").(string))
			if slices.Contains(geoReplicationLocations, location) {
				return errors.New("the `georeplications` list cannot contain the location where the Container Registry exists")
			}

			quarantinePolicyEnabled := d.Get("quarantine_policy_enabled").(bool)
			if quarantinePolicyEnabled && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return errors.New("an ACR quarantine policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please unset quarantine_policy_enabled")
			}

			retentionPolicyEnabled, ok := d.GetOk("retention_policy_in_days")
			if ok && retentionPolicyEnabled.(int) > 0 && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return errors.New("an ACR retention policy can only be applied when using the Premium Sku. If you are downgrading from a Premium SKU please unset `retention_policy_in_days`")
			}

			exportPolicyEnabled := d.Get("export_policy_enabled").(bool)
			if !exportPolicyEnabled {
				if !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
					return errors.New("an ACR export policy can only be disabled when using the Premium Sku. If you are downgrading from a Premium SKU please unset `export_policy_enabled` or set `export_policy_enabled = true`")
				}
				if d.Get("public_network_access_enabled").(bool) {
					return errors.New("to disable export of artifacts, `public_network_access_enabled` must also be `false`")
				}
			}

			encryptionEnabled, ok := d.GetOk("encryption")
			if ok && len(encryptionEnabled.([]any)) > 0 && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return errors.New("an ACR encryption can only be applied when using the Premium Sku")
			}

			// zone redundancy is only available for Premium Sku.
			zoneRedundancyEnabled, ok := d.GetOk("zone_redundancy_enabled")
			if ok && zoneRedundancyEnabled.(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("ACR zone redundancy can only be applied when using the Premium Sku")
			}
			for _, loc := range geoReplications {
				loc := loc.(map[string]any)
				zoneRedundancyEnabled, ok := loc["zone_redundancy_enabled"]
				if ok && zoneRedundancyEnabled.(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
					return fmt.Errorf("ACR zone redundancy can only be applied when using the Premium Sku")
				}
			}

			// anonymous pull is only available for Standard/Premium Sku.
			if d.Get("anonymous_pull_enabled").(bool) && (!strings.EqualFold(sku, string(registries.SkuNameStandard)) && !strings.EqualFold(sku, string(registries.SkuNamePremium))) {
				return fmt.Errorf("`anonymous_pull_enabled` can only be applied when using the Standard/Premium Sku")
			}

			// data endpoint is only available for Premium Sku.
			if d.Get("data_endpoint_enabled").(bool) && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
				return fmt.Errorf("`data_endpoint_enabled` can only be applied when using the Premium Sku")
			}

			return nil
		}),
	}
}
