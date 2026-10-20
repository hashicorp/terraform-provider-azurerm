// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace", id.ID())
		}
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	location := location.Normalize(d.Get("location").(string))
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	t := d.Get("tags").(map[string]any)
	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	publicNetworkEnabled := namespaces.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkEnabled = namespaces.PublicNetworkAccessDisabled
	}

	disableLocalAuth := !d.Get("local_authentication_enabled").(bool)

	parameters := namespaces.EHNamespace{
		Location: &location,
		Sku: &namespaces.Sku{
			Name: namespaces.SkuName(sku),
			Tier: func() *namespaces.SkuTier {
				return pointer.ToEnum[namespaces.SkuTier](sku)
			}(),
			Capacity: pointer.To(int64(capacity)),
		},
		Identity: identity,
		Properties: &namespaces.EHNamespaceProperties{
			IsAutoInflateEnabled: pointer.To(autoInflateEnabled),
			DisableLocalAuth:     pointer.To(disableLocalAuth),
			PublicNetworkAccess:  &publicNetworkEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v := d.Get("dedicated_cluster_id").(string); v != "" {
		parameters.Properties.ClusterArmId = pointer.To(v)
	}

	if tlsValue := d.Get("minimum_tls_version").(string); tlsValue != "" {
		parameters.Properties.MinimumTlsVersion = pointer.ToEnum[namespaces.TlsVersion](tlsValue)
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		parameters.Properties.MaximumThroughputUnits = pointer.To(int64(v.(int)))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	ruleSets, hasRuleSets := d.GetOk("network_rulesets")
	if hasRuleSets {
		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name == namespaces.SkuNameBasic {
			return fmt.Errorf("network_rulesets cannot be used when the SKU is basic")
		}

		ruleSetProperties, err := expandEventHubNamespaceNetworkRuleset(ruleSets.([]any))
		if err != nil {
			return err
		}
		rulesets := networkrulesets.NetworkRuleSet{
			Properties: ruleSetProperties,
		}

		if !strings.EqualFold(string(*rulesets.Properties.PublicNetworkAccess), string(*parameters.Properties.PublicNetworkAccess)) {
			return fmt.Errorf("the value of public network access of namespace should be the same as of the network rulesets")
		}

		ruleSetsClient := meta.(*clients.Client).Eventhub.NetworkRuleSetsClient
		namespaceId := networkrulesets.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
		if _, err := ruleSetsClient.NamespacesCreateOrUpdateNetworkRuleSet(ctx, namespaceId, rulesets); err != nil {
			return fmt.Errorf("setting network ruleset properties for %s: %+v", id, err)
		}
	}

	return resourceEventHubNamespaceRead(d, meta)
}

func expandEventHubNamespaceNetworkRuleset(input []any) (*networkrulesets.NetworkRuleSetProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}

	block := input[0].(map[string]any)

	publicNetworkAccess := networkrulesets.PublicNetworkAccessFlagEnabled
	if !block["public_network_access_enabled"].(bool) {
		publicNetworkAccess = networkrulesets.PublicNetworkAccessFlagDisabled
	}

	defaultAction := networkrulesets.DefaultAction(block["default_action"].(string))
	ruleset := networkrulesets.NetworkRuleSetProperties{
		DefaultAction:       &defaultAction,
		PublicNetworkAccess: &publicNetworkAccess,
	}

	if v, ok := block["trusted_service_access_enabled"]; ok {
		ruleset.TrustedServiceAccessEnabled = pointer.To(v.(bool))
	}

	if v, ok := block["virtual_network_rule"]; ok {
		value := v.(*pluginsdk.Set).List()
		if len(value) > 0 {
			var rules []networkrulesets.NWRuleSetVirtualNetworkRules
			for _, r := range value {
				rblock := r.(map[string]any)
				rules = append(rules, networkrulesets.NWRuleSetVirtualNetworkRules{
					Subnet: &networkrulesets.Subnet{
						Id: pointer.To(rblock["subnet_id"].(string)),
					},
					IgnoreMissingVnetServiceEndpoint: pointer.To(rblock["ignore_missing_virtual_network_service_endpoint"].(bool)),
				})
			}

			ruleset.VirtualNetworkRules = &rules
		}
	}

	if v, ok := block["ip_rule"].([]any); ok {
		if len(v) > 0 {
			var rules []networkrulesets.NWRuleSetIPRules
			for _, r := range v {
				rblock := r.(map[string]any)
				rules = append(rules, networkrulesets.NWRuleSetIPRules{
					IPMask: pointer.To(rblock["ip_mask"].(string)),
					Action: func() *networkrulesets.NetworkRuleIPAction {
						return pointer.ToEnum[networkrulesets.NetworkRuleIPAction](rblock["action"].(string))
					}(),
				})
			}

			ruleset.IPRules = &rules
		}
	}

	// the API silently keeps the effective `default_action` as `Allow` if `Deny` is sent with no `ip_rule` or
	// `virtual_network_rule` defined, and returns no error - producing permanent plan drift on the next read
	if defaultAction == networkrulesets.DefaultActionDeny && ruleset.VirtualNetworkRules == nil && ruleset.IPRules == nil {
		return nil, errors.New("the `default_action` of `network_rulesets` can only be set to `Deny` when at least one `ip_rule` or `virtual_network_rule` block is specified")
	}

	return &ruleset, nil
}
