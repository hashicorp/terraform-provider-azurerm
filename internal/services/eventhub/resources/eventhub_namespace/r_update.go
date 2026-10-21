// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	location := location.Normalize(d.Get("location").(string))
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	t := d.Get("tags").(map[string]any)
	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)

	publicNetworkEnabled := namespaces.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkEnabled = namespaces.PublicNetworkAccessDisabled
	}

	disableLocalAuth := !d.Get("local_authentication_enabled").(bool)

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

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

	if d.HasChange("maximum_throughput_units") {
		parameters.Properties.MaximumThroughputUnits = pointer.To(int64(d.Get("maximum_throughput_units").(int)))
	}

	// @favoretti: if we are downgrading from Standard to Basic SKU and namespace had both autoInflate enabled and
	// maximumThroughputUnits set - we need to force throughput units back to 0, otherwise downgrade fails
	//
	// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/10244
	//
	if *parameters.Sku.Tier == namespaces.SkuTierBasic && !autoInflateEnabled {
		parameters.Properties.MaximumThroughputUnits = pointer.To(int64(0))
	}

	if _, err = client.Update(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	deadline, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"Activating", "ActivatingIdentity", "Updating", "Pending"},
		Target:       []string{"Succeeded"},
		Refresh:      eventHubNamespaceProvisioningStateRefreshFunc(ctx, client, id),
		Timeout:      time.Until(deadline),
		PollInterval: 10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be updated: %+v", id, err)
	}

	if d.HasChange("network_rulesets") {
		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name == namespaces.SkuNameBasic {
			return fmt.Errorf("network_rulesets cannot be used when the SKU is basic")
		}

		ruleSets := d.Get("network_rulesets")
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
