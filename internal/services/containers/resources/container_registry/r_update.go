// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.Registries
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registries.ParseRegistryID(d.Id())
	if err != nil {
		return err
	}

	sku := d.Get("sku").(string)
	skuChange := d.HasChange("sku")
	isBasicSku := strings.EqualFold(sku, string(registries.SkuNameBasic))
	isPremiumSku := strings.EqualFold(sku, string(registries.SkuNamePremium))
	isStandardSku := strings.EqualFold(sku, string(registries.SkuNameStandard))

	oldReplicationsRaw, newReplicationsRaw := d.GetChange("georeplications")
	hasGeoReplicationsChanges := d.HasChange("georeplications")
	oldReplications := oldReplicationsRaw.([]any)
	newReplications := newReplicationsRaw.([]any)

	// handle upgrade to Premium SKU first
	if skuChange && isPremiumSku {
		if err := applyContainerRegistrySku(d, meta, sku, *id); err != nil {
			return fmt.Errorf("applying sku %q for %s: %+v", sku, id, err)
		}
	}

	payload := &registries.RegistryUpdateParameters{
		Properties: &registries.RegistryPropertiesUpdateParameters{},
	}

	if d.HasChange("network_rule_set") {
		networkRuleSet := expandNetworkRuleSet(d.Get("network_rule_set").([]any))
		if networkRuleSet != nil && isBasicSku {
			return fmt.Errorf("`network_rule_set` can only be specified for a Premium Sku. If you are reverting from a Premium to Basic SKU please set network_rule_set = []")
		}

		payload.Properties.NetworkRuleSet = networkRuleSet
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := registries.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			if !isPremiumSku {
				return errors.New("`public_network_access_enabled` can only be disabled for a Premium Sku")
			}

			publicNetworkAccess = registries.PublicNetworkAccessDisabled
		}

		payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		payload.Identity = identity
	}

	policyKeys := []string{
		"quarantine_policy_enabled",
		"export_policy_enabled",
		"azuread_authentication_as_arm_policy_enabled",
	}

	policyKeys = append(policyKeys, "retention_policy_in_days")

	if d.HasChanges(policyKeys...) {
		payload.Properties.Policies = &registries.Policies{}
	}

	if d.HasChange("retention_policy_in_days") {
		payload.Properties.Policies.RetentionPolicy = &registries.RetentionPolicy{
			Status: pointer.To(registries.PolicyStatusDisabled),
		}

		if v := d.Get("retention_policy_in_days").(int); v != 0 {
			payload.Properties.Policies.RetentionPolicy = &registries.RetentionPolicy{
				Status: pointer.To(registries.PolicyStatusEnabled),
				Days:   pointer.To(int64(v)),
			}
		}
	}

	if d.HasChange("quarantine_policy_enabled") {
		payload.Properties.Policies.QuarantinePolicy = &registries.QuarantinePolicy{
			Status: pointer.To(registries.PolicyStatusDisabled),
		}

		if v := d.Get("quarantine_policy_enabled").(bool); v {
			payload.Properties.Policies.QuarantinePolicy = &registries.QuarantinePolicy{
				Status: pointer.To(registries.PolicyStatusEnabled),
			}
		}
	}

	if d.HasChange("export_policy_enabled") {
		payload.Properties.Policies.ExportPolicy = &registries.ExportPolicy{
			Status: pointer.To(registries.ExportPolicyStatusDisabled),
		}

		if v := d.Get("export_policy_enabled").(bool); v {
			payload.Properties.Policies.ExportPolicy = &registries.ExportPolicy{
				Status: pointer.To(registries.ExportPolicyStatusEnabled),
			}
		}
	}

	if d.HasChange("azuread_authentication_as_arm_policy_enabled") {
		payload.Properties.Policies.AzureADAuthenticationAsArmPolicy = expandAadAuthAsArmPolicy(d.Get("azuread_authentication_as_arm_policy_enabled").(bool))
	}

	if d.HasChange("admin_enabled") {
		payload.Properties.AdminUserEnabled = pointer.To(d.Get("admin_enabled").(bool))
	}

	if d.HasChange("encryption") {
		payload.Properties.Encryption = expandEncryption(d.Get("encryption").([]any))
	}

	if d.HasChange("anonymous_pull_enabled") {
		payload.Properties.AnonymousPullEnabled = pointer.To(d.Get("anonymous_pull_enabled").(bool))
	}

	if d.HasChange("data_endpoint_enabled") {
		payload.Properties.DataEndpointEnabled = pointer.To(d.Get("data_endpoint_enabled").(bool))
	}

	if d.HasChange("network_rule_bypass_option") {
		payload.Properties.NetworkRuleBypassOptions = pointer.ToEnum[registries.NetworkRuleBypassOptions](d.Get("network_rule_bypass_option").(string))
	}

	if d.HasChange("role_assignment_mode") {
		payload.Properties.RoleAssignmentMode = pointer.ToEnum[registries.RoleAssignmentMode](d.Get("role_assignment_mode").(string))
	}

	if d.HasChange("network_rule_bypass_for_tasks_enabled") {
		payload.Properties.NetworkRuleBypassAllowedForTasks = pointer.To(d.Get("network_rule_bypass_for_tasks_enabled").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	// geo replication is only supported by Premium Sku
	if len(newReplications) > 0 && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
		return errors.New("an ACR geo-replication can only be applied when using the Premium Sku")
	}

	if hasGeoReplicationsChanges {
		if err := applyGeoReplicationLocations(ctx, meta, *id, expandReplications(oldReplications), expandReplications(newReplications)); err != nil {
			return fmt.Errorf("applying geo replications for %s: %+v", id, err)
		}
	}

	if err := client.UpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// downgrade to Basic or Standard SKU
	if skuChange && (isBasicSku || isStandardSku) {
		if err := applyContainerRegistrySku(d, meta, sku, *id); err != nil {
			return fmt.Errorf("applying sku %q for %s: %+v", sku, id, err)
		}
	}

	d.SetId(id.ID())

	return resourceContainerRegistryRead(d, meta)
}
