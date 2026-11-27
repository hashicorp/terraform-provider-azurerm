// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).NamespacesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	location := location.Normalize(d.Get("location").(string))
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]any)

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace", id.ID())
		}
	}

	identity, err := expandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	publicNetworkEnabled := namespaces.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkEnabled = namespaces.PublicNetworkAccessDisabled
	}

	parameters := namespaces.SBNamespace{
		Location: location,
		Identity: identity,
		Sku: &namespaces.SBSku{
			Name: namespaces.SkuName(sku),
			Tier: pointer.ToEnum[namespaces.SkuTier](sku),
		},
		Properties: &namespaces.SBNamespaceProperties{
			Encryption:          expandServiceBusNamespaceEncryption(d.Get("customer_managed_key").([]any)),
			DisableLocalAuth:    pointer.To(!d.Get("local_auth_enabled").(bool)),
			PublicNetworkAccess: &publicNetworkEnabled,
		},
		Tags: expandTags(t),
	}

	if tlsValue := d.Get("minimum_tls_version").(string); tlsValue != "" {
		parameters.Properties.MinimumTlsVersion = pointer.ToEnum[namespaces.TlsVersion](tlsValue)
	}

	if capacity := d.Get("capacity"); capacity != nil {
		if !strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) > 0 {
			return fmt.Errorf("service bus SKU %q only supports `capacity` of 0", sku)
		}
		if strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && capacity.(int) == 0 {
			return fmt.Errorf("service bus SKU %q only supports `capacity` of 1, 2, 4, 8 or 16", sku)
		}
		parameters.Sku.Capacity = pointer.To(int64(capacity.(int)))
	}

	if premiumMessagingUnit := d.Get("premium_messaging_partitions"); premiumMessagingUnit != nil {
		if !strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && premiumMessagingUnit.(int) > 0 {
			return fmt.Errorf("premium messaging partition is not supported by service Bus SKU %q and it can only be set to 0", sku)
		}
		if strings.EqualFold(sku, string(namespaces.SkuNamePremium)) && premiumMessagingUnit.(int) == 0 {
			return fmt.Errorf("service bus SKU %q only supports `premium_messaging_partitions` of 1, 2, 4", sku)
		}
		parameters.Properties.PremiumMessagingPartitions = pointer.To(int64(premiumMessagingUnit.(int)))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	if err = createNetworkRuleSetForNamespace(ctx, client, id, d.Get("network_rule_set").([]any)); err != nil {
		return err
	}

	return resourceServiceBusNamespaceRead(d, meta)
}
