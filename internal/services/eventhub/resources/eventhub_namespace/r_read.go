// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	authorizationKeysClient := meta.(*clients.Client).Eventhub.NamespaceAuthorizationRulesClient
	ruleSetsClient := meta.(*clients.Client).Eventhub.NetworkRuleSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(sku.Name))
			d.Set("capacity", sku.Capacity)
		}

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
			d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
			d.Set("dedicated_cluster_id", props.ClusterArmId)

			d.Set("local_authentication_enabled", !pointer.From(props.DisableLocalAuth))

			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == namespaces.PublicNetworkAccessDisabled {
				publicNetworkAccess = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)
			d.Set("minimum_tls_version", pointer.FromEnum(props.MinimumTlsVersion))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	namespaceId := networkrulesets.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	ruleset, err := ruleSetsClient.NamespacesGetNetworkRuleSet(ctx, namespaceId)
	if err != nil {
		return fmt.Errorf("retrieving Network Rule Sets for %s: %+v", *id, err)
	}

	networkRuleSets, err := flattenEventHubNamespaceNetworkRuleset(ruleset)
	if err != nil {
		return fmt.Errorf("flattening `network_rule` for %s: %+v", id, err)
	}
	if err := d.Set("network_rulesets", networkRuleSets); err != nil {
		return fmt.Errorf("setting `network_ruleset` for Evenhub Namespace %s: %v", id.NamespaceName, err)
	}

	authorizationRuleId := authorizationrulesnamespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, eventHubNamespaceDefaultAuthorizationRule)
	keys, err := authorizationKeysClient.NamespacesListKeys(ctx, authorizationRuleId)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for EventHub Namespace %q: %+v", id.NamespaceName, err)
	}

	if model := keys.Model; model != nil {
		d.Set("default_primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("default_secondary_connection_string_alias", model.AliasSecondaryConnectionString)
		d.Set("default_primary_connection_string", model.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", model.SecondaryConnectionString)
		d.Set("default_primary_key", model.PrimaryKey)
		d.Set("default_secondary_key", model.SecondaryKey)
	}

	return nil
}

func flattenEventHubNamespaceNetworkRuleset(ruleset networkrulesets.NamespacesGetNetworkRuleSetOperationResponse) ([]any, error) {
	if ruleset.Model == nil || ruleset.Model.Properties == nil {
		return []any{}, nil
	}

	vnetBlocks := make([]any, 0)
	if vnetRules := ruleset.Model.Properties.VirtualNetworkRules; vnetRules != nil {
		for _, vnetRule := range *vnetRules {
			block := make(map[string]any)

			if s := vnetRule.Subnet; s != nil {
				if v := s.Id; v != nil {
					// the API returns the subnet ID's resource group name in lowercase
					// https://github.com/Azure/azure-sdk-for-go/issues/5855
					// for some reason the DiffSuppressFunc for `subnet_id` isn't working as intended, so we'll also flatten the id insensitively
					subnetId, err := commonids.ParseSubnetIDInsensitively(*v)
					if err != nil {
						return nil, fmt.Errorf("parsing `subnet_id`: %+v", err)
					}
					block["subnet_id"] = subnetId.ID()
				}
			}

			if v := vnetRule.IgnoreMissingVnetServiceEndpoint; v != nil {
				block["ignore_missing_virtual_network_service_endpoint"] = *v
			}

			vnetBlocks = append(vnetBlocks, block)
		}
	}
	ipBlocks := make([]any, 0)
	if ipRules := ruleset.Model.Properties.IPRules; ipRules != nil {
		for _, ipRule := range *ipRules {
			block := make(map[string]any)

			action := ""
			if ipRule.Action != nil {
				action = string(*ipRule.Action)
			}

			block["action"] = action

			if v := ipRule.IPMask; v != nil {
				block["ip_mask"] = *v
			}

			ipBlocks = append(ipBlocks, block)
		}
	}

	// TODO: fix this

	publicNetworkAccess := ruleset.Model.Properties.PublicNetworkAccess == nil || *ruleset.Model.Properties.PublicNetworkAccess != networkrulesets.PublicNetworkAccessFlagDisabled

	return []any{map[string]any{
		"default_action":                 string(*ruleset.Model.Properties.DefaultAction),
		"public_network_access_enabled":  publicNetworkAccess,
		"virtual_network_rule":           vnetBlocks,
		"ip_rule":                        ipBlocks,
		"trusted_service_access_enabled": ruleset.Model.Properties.TrustedServiceAccessEnabled,
	}}, nil
}

// The resource id of subnet_id that's being returned by API is always lower case &
// the default caseDiff suppress func is not working in TypeSet
func resourceVnetRuleHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		if v, ok := m["subnet_id"]; ok {
			fmt.Fprintf(&buf, "%s-", strings.ToLower(v.(string)))
		}
		if v, ok := m["ignore_missing_virtual_network_service_endpoint"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
	}
	return pluginsdk.HashString(buf.String())
}

func eventhubTLSVersionDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) (err error) {
	old, new := d.GetChange("minimum_tls_version")
	if old != "" && new == "" {
		err = fmt.Errorf("`minimum_tls_version` has been set before, please set a valid value for this property ")
	}
	return
}

func eventHubNamespaceProvisioningStateRefreshFunc(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)

		provisioningState := "Pending"
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, provisioningState, nil
			}
			return nil, "Error", fmt.Errorf("polling for the provisioning state of %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties != nil && res.Model.Properties.ProvisioningState != nil {
			provisioningState = *res.Model.Properties.ProvisioningState
		}

		return res, provisioningState, nil
	}
}
