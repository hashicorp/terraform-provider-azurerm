// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespacesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceServiceBusNamespaceFlatten(ctx context.Context, d *pluginsdk.ResourceData, id *namespaces.NamespaceId, model *namespaces.SBNamespace, client *namespaces.NamespacesClient, namespaceAuthClient *namespacesauthorizationrule.NamespacesAuthorizationRuleClient) error {
	d.Set("name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		d.Set("tags", flattenTags(model.Tags))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if sku := model.Sku; sku != nil {
			skuName := ""
			// the Azure API is inconsistent here, so rewrite this into the casing we expect
			for _, v := range namespaces.PossibleValuesForSkuName() {
				if strings.EqualFold(v, string(sku.Name)) {
					skuName = v
				}
			}
			d.Set("sku", skuName)
			d.Set("capacity", sku.Capacity)

			if props := model.Properties; props != nil {
				d.Set("premium_messaging_partitions", props.PremiumMessagingPartitions)

				if customerManagedKey, err := flattenServiceBusNamespaceEncryption(props.Encryption); err == nil {
					d.Set("customer_managed_key", customerManagedKey)
				}
				localAuthEnabled := true
				if props.DisableLocalAuth != nil {
					localAuthEnabled = !*props.DisableLocalAuth
				}
				d.Set("local_auth_enabled", localAuthEnabled)

				publicNetworkAccess := true
				if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == namespaces.PublicNetworkAccessDisabled {
					publicNetworkAccess = false
				}
				d.Set("public_network_access_enabled", publicNetworkAccess)

				if props.MinimumTlsVersion != nil {
					d.Set("minimum_tls_version", pointer.FromEnum(props.MinimumTlsVersion))
				}

				d.Set("endpoint", props.ServiceBusEndpoint)
			}
		}
	}

	authRuleId := namespacesauthorizationrule.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, serviceBusNamespaceDefaultAuthorizationRule)

	keys, err := namespaceAuthClient.NamespacesListKeys(ctx, authRuleId)
	if err != nil {
		log.Printf("[WARN] listing default keys for %s: %+v", id, err)
	} else {
		if keysModel := keys.Model; keysModel != nil {
			d.Set("default_primary_connection_string", keysModel.PrimaryConnectionString)
			d.Set("default_secondary_connection_string", keysModel.SecondaryConnectionString)
			d.Set("default_primary_key", keysModel.PrimaryKey)
			d.Set("default_secondary_key", keysModel.SecondaryKey)
		}
	}

	networkRuleSet, err := client.GetNetworkRuleSet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving network rule set %s: %+v", *id, err)
	}

	if nrsModel := networkRuleSet.Model; nrsModel != nil {
		if props := nrsModel.Properties; props != nil {
			d.Set("network_rule_set", flattenServiceBusNamespaceNetworkRuleSet(*props))
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenServiceBusNamespaceEncryption(encryption *namespaces.Encryption) ([]any, error) {
	if encryption == nil {
		return []any{}, nil
	}

	var keyId string
	var identityId string
	if keyVaultProperties := encryption.KeyVaultProperties; keyVaultProperties != nil && len(*keyVaultProperties) != 0 {
		props := (*keyVaultProperties)[0]
		keyVaultKeyId, err := keyvault.NewNestedItemID(pointer.From(props.KeyVaultUri), keyvault.NestedItemTypeKey, pointer.From(props.KeyName), pointer.From(props.KeyVersion))
		if err != nil {
			return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		keyId = keyVaultKeyId.ID()
		if props.Identity != nil && props.Identity.UserAssignedIdentity != nil {
			sbnUaiId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*props.Identity.UserAssignedIdentity)
			if err != nil {
				return nil, err
			}
			identityId = sbnUaiId.ID()
		}
	}

	return []any{
		map[string]any{
			"infrastructure_encryption_enabled": encryption.RequireInfrastructureEncryption,
			"key_vault_key_id":                  keyId,
			"identity_id":                       identityId,
		},
	}, nil
}

func flattenServiceBusNamespaceNetworkRuleSet(networkRuleSet namespaces.NetworkRuleSetProperties) []any {
	defaultAction := ""
	if v := networkRuleSet.DefaultAction; v != nil {
		defaultAction = string(*v)
	}
	publicNetworkAccess := namespaces.PublicNetworkAccessFlagEnabled
	if v := networkRuleSet.PublicNetworkAccess; v != nil {
		publicNetworkAccess = *v
	}

	networkRules := flattenServiceBusNamespaceVirtualNetworkRules(networkRuleSet.VirtualNetworkRules)
	ipRules := flattenServiceBusNamespaceIPRules(networkRuleSet.IPRules)

	return []any{map[string]any{
		"default_action":                defaultAction,
		"trusted_services_allowed":      pointer.From(networkRuleSet.TrustedServiceAccessEnabled),
		"public_network_access_enabled": publicNetworkAccess == namespaces.PublicNetworkAccessFlagEnabled,
		"network_rules":                 pluginsdk.NewSet(networkRuleHash, networkRules),
		"ip_rules":                      ipRules,
	}}
}

func flattenServiceBusNamespaceVirtualNetworkRules(input *[]namespaces.NWRuleSetVirtualNetworkRules) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0, len(*input))
	for _, v := range *input {
		subnetId := ""
		if v.Subnet != nil && v.Subnet.Id != "" {
			subnetId = v.Subnet.Id
		}

		result = append(result, map[string]any{
			"subnet_id":                            subnetId,
			"ignore_missing_vnet_service_endpoint": pointer.From(v.IgnoreMissingVnetServiceEndpoint),
		})
	}

	return result
}

func flattenServiceBusNamespaceIPRules(input *[]namespaces.NWRuleSetIPRules) []any {
	if input == nil || len(*input) == 0 {
		return []any{}
	}

	result := make([]any, 0, len(*input))
	for _, v := range *input {
		if v.IPMask != nil {
			result = append(result, *v.IPMask)
		}
	}

	return result
}
