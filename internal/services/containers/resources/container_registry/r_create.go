// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/replications"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.Registries
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := registries.NewRegistryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry", id.ID())
		}
	}

	sku := d.Get("sku").(string)

	networkRuleSet := expandNetworkRuleSet(d.Get("network_rule_set").([]any))
	if networkRuleSet != nil && !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
		return fmt.Errorf("`network_rule_set` can only be specified for a Premium Sku. If you are reverting from a Premium to Basic SKU please set network_rule_set = []")
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	publicNetworkAccess := registries.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		if !strings.EqualFold(sku, string(registries.SkuNamePremium)) {
			return errors.New("`public_network_access_enabled` can only be disabled for a Premium Sku")
		}

		publicNetworkAccess = registries.PublicNetworkAccessDisabled
	}

	zoneRedundancy := registries.ZoneRedundancyDisabled
	if d.Get("zone_redundancy_enabled").(bool) {
		zoneRedundancy = registries.ZoneRedundancyEnabled
	}

	retentionPolicy := &registries.RetentionPolicy{}
	if v, ok := d.GetOk("retention_policy_in_days"); ok && v.(int) > 0 {
		retentionPolicy.Days = pointer.To(int64(v.(int)))
		retentionPolicy.Status = pointer.To(registries.PolicyStatusEnabled)
	}

	parameters := registries.Registry{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: registries.Sku{
			Name: registries.SkuName(sku),
			Tier: pointer.ToEnum[registries.SkuTier](sku),
		},
		Identity: identity,
		Properties: &registries.RegistryProperties{
			AdminUserEnabled: pointer.To(d.Get("admin_enabled").(bool)),
			Encryption:       expandEncryption(d.Get("encryption").([]any)),
			NetworkRuleSet:   networkRuleSet,
			Policies: &registries.Policies{
				QuarantinePolicy:                 expandQuarantinePolicy(d.Get("quarantine_policy_enabled").(bool)),
				RetentionPolicy:                  retentionPolicy,
				ExportPolicy:                     expandExportPolicy(d.Get("export_policy_enabled").(bool)),
				AzureADAuthenticationAsArmPolicy: expandAadAuthAsArmPolicy(d.Get("azuread_authentication_as_arm_policy_enabled").(bool)),
			},
			PublicNetworkAccess:              &publicNetworkAccess,
			ZoneRedundancy:                   &zoneRedundancy,
			AnonymousPullEnabled:             pointer.To(d.Get("anonymous_pull_enabled").(bool)),
			DataEndpointEnabled:              pointer.To(d.Get("data_endpoint_enabled").(bool)),
			NetworkRuleBypassOptions:         pointer.ToEnum[registries.NetworkRuleBypassOptions](d.Get("network_rule_bypass_option").(string)),
			RoleAssignmentMode:               pointer.ToEnum[registries.RoleAssignmentMode](d.Get("role_assignment_mode").(string)),
			NetworkRuleBypassAllowedForTasks: pointer.To(d.Get("network_rule_bypass_for_tasks_enabled").(bool)),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// the ACR is being created so no previous geo-replication locations
	var oldGeoReplicationLocations, newGeoReplicationLocations []replications.Replication
	newGeoReplicationLocations = expandReplications(d.Get("georeplications").([]any))
	// geo replications have been specified
	if len(newGeoReplicationLocations) > 0 {
		if err = applyGeoReplicationLocations(ctx, meta, id, oldGeoReplicationLocations, newGeoReplicationLocations); err != nil {
			return fmt.Errorf("applying geo replications for %s: %+v", id, err)
		}
	}

	return resourceContainerRegistryRead(d, meta)
}
