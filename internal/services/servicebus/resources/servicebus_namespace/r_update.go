// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).NamespacesClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}

	payload := existing.Model

	if d.HasChange("identity") {
		identity, err := expandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = identity
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkEnabled := namespaces.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkEnabled = namespaces.PublicNetworkAccessDisabled
		}
		payload.Properties.PublicNetworkAccess = &publicNetworkEnabled
	}

	if d.HasChange("sku") {
		sku := d.Get("sku").(string)
		payload.Sku = &namespaces.SBSku{
			Name: namespaces.SkuName(sku),
			Tier: pointer.ToEnum[namespaces.SkuTier](sku),
		}
	}

	if d.HasChange("customer_managed_key") {
		payload.Properties.Encryption = expandServiceBusNamespaceEncryption(d.Get("customer_managed_key").([]any))
	}

	if d.HasChange("local_auth_enabled") {
		payload.Properties.DisableLocalAuth = pointer.To(!d.Get("local_auth_enabled").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = expandTags(d.Get("tags").(map[string]any))
	}

	if d.HasChange("minimum_tls_version") {
		payload.Properties.MinimumTlsVersion = pointer.ToEnum[namespaces.TlsVersion](d.Get("minimum_tls_version").(string))
	}

	if d.HasChange("capacity") {
		sku := d.Get("sku").(string)
		if capacity := d.Get("capacity"); capacity != nil {
			if !strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) > 0 {
				return fmt.Errorf("service bus SKU %q only supports `capacity` of 0", sku)
			}
			if strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) == 0 {
				return fmt.Errorf("service bus SKU %q only supports `capacity` of 1, 2, 4, 8 or 16", sku)
			}
			payload.Sku.Capacity = pointer.To(int64(capacity.(int)))
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, id); err != nil {
		return err
	}

	if d.HasChange("network_rule_set") {
		oldNetworkRuleSet, newNetworkRuleSet := d.GetChange("network_rule_set")
		// if the network rule set has been removed from config, reset it instead as there is no way to remove a rule set
		if len(oldNetworkRuleSet.([]any)) == 1 && len(newNetworkRuleSet.([]any)) == 0 {
			log.Printf("[DEBUG] Resetting Network Rule Set associated with %s..", id)
			if err = resetNetworkRuleSetForNamespace(ctx, client, *id); err != nil {
				return err
			}
		} else {
			if err = createNetworkRuleSetForNamespace(ctx, client, *id, newNetworkRuleSet.([]any)); err != nil {
				return err
			}
		}
	}

	return resourceServiceBusNamespaceRead(d, meta)
}
