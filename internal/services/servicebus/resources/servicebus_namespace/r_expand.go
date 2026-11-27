// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandServiceBusNamespaceEncryption(input []any) *namespaces.Encryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	keyId, _ := keyvault.ParseNestedItemID(v["key_vault_key_id"].(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)

	encryption := namespaces.Encryption{
		KeySource:                       pointer.To(namespaces.KeySourceMicrosoftPointKeyVault),
		RequireInfrastructureEncryption: pointer.To(v["infrastructure_encryption_enabled"].(bool)),
	}

	encryption.KeyVaultProperties = &[]namespaces.KeyVaultProperties{
		{
			KeyName:     pointer.To(keyId.Name),
			KeyVersion:  pointer.To(keyId.Version),
			KeyVaultUri: pointer.To(keyId.KeyVaultBaseURL),
			Identity: &namespaces.UserAssignedIdentityProperties{
				UserAssignedIdentity: pointer.To(v["identity_id"].(string)),
			},
		},
	}

	return &encryption
}

func expandSystemAndUserAssignedMap(input []any) (*identity.SystemAndUserAssignedMap, error) {
	identityType := identity.TypeNone
	identityIds := make(map[string]identity.UserAssignedIdentityDetails, 0)

	if len(input) > 0 {
		raw := input[0].(map[string]any)
		typeRaw := raw["type"].(string)
		if typeRaw == string(identity.TypeSystemAssigned) {
			identityType = identity.TypeSystemAssigned
		}
		if typeRaw == string(identity.TypeSystemAssignedUserAssigned) {
			identityType = identity.TypeSystemAssignedUserAssigned
		}
		if typeRaw == string(identity.TypeUserAssigned) {
			identityType = identity.TypeUserAssigned
		}

		identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
		for _, v := range identityIdsRaw {
			identityIds[v.(string)] = identity.UserAssignedIdentityDetails{
				// intentionally empty since the expand shouldn't send these values
			}
		}
	}

	if len(identityIds) > 0 && (identityType != identity.TypeSystemAssignedUserAssigned && identityType != identity.TypeUserAssigned) {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q or %q", string(identity.TypeSystemAssignedUserAssigned), string(identity.TypeUserAssigned))
	}

	if len(identityIds) == 0 {
		return &identity.SystemAndUserAssignedMap{
			Type: identityType,
		}, nil
	}

	return &identity.SystemAndUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}

func expandServiceBusNamespaceVirtualNetworkRules(input []any) *[]namespaces.NWRuleSetVirtualNetworkRules {
	if len(input) == 0 {
		return nil
	}

	result := make([]namespaces.NWRuleSetVirtualNetworkRules, 0)
	for _, v := range input {
		raw := v.(map[string]any)
		result = append(result, namespaces.NWRuleSetVirtualNetworkRules{
			Subnet: &namespaces.Subnet{
				Id: raw["subnet_id"].(string),
			},
			IgnoreMissingVnetServiceEndpoint: pointer.To(raw["ignore_missing_vnet_service_endpoint"].(bool)),
		})
	}

	return &result
}

func expandServiceBusNamespaceIPRules(input []any) *[]namespaces.NWRuleSetIPRules {
	if len(input) == 0 {
		return nil
	}

	result := make([]namespaces.NWRuleSetIPRules, 0, len(input))
	for _, v := range input {
		result = append(result, namespaces.NWRuleSetIPRules{
			IPMask: pointer.To(v.(string)),
			Action: pointer.To(namespaces.NetworkRuleIPActionAllow),
		})
	}

	return &result
}

func createNetworkRuleSetForNamespace(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId, input []any) error {
	if len(input) < 1 || input[0] == nil {
		return nil
	}

	item := input[0].(map[string]any)

	defaultAction := namespaces.DefaultAction(item["default_action"].(string))
	vnetRule := expandServiceBusNamespaceVirtualNetworkRules(item["network_rules"].(*pluginsdk.Set).List())
	ipRule := expandServiceBusNamespaceIPRules(item["ip_rules"].(*pluginsdk.Set).List())
	publicNetworkAcc := "Disabled"
	if item["public_network_access_enabled"].(bool) {
		publicNetworkAcc = "Enabled"
	}

	// API doesn't accept "Deny" to be set for "default_action" if no "ip_rules" or "network_rules" is defined and returns no error message to the user
	if defaultAction == namespaces.DefaultActionDeny && vnetRule == nil && ipRule == nil {
		return errors.New(" The `default_action` of `network_rule_set` can only be set to `Allow` if no `ip_rules` or `network_rules` is set")
	}

	parameters := namespaces.NetworkRuleSet{
		Properties: &namespaces.NetworkRuleSetProperties{
			DefaultAction:               &defaultAction,
			VirtualNetworkRules:         vnetRule,
			IPRules:                     ipRule,
			PublicNetworkAccess:         pointer.ToEnum[namespaces.PublicNetworkAccessFlag](publicNetworkAcc),
			TrustedServiceAccessEnabled: pointer.To(item["trusted_services_allowed"].(bool)),
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	return nil
}
