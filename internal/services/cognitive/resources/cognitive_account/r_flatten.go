// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceCognitiveAccountFlatten(ctx context.Context, client *cognitiveservicesaccounts.CognitiveServicesAccountsClient, d *pluginsdk.ResourceData, id *cognitiveservicesaccounts.AccountId, account *cognitiveservicesaccounts.Account, fetchAccessKeys bool) error {
	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if account != nil {
		d.Set("kind", account.Kind)

		d.Set("location", location.NormalizeNilable(account.Location))
		if sku := account.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}

		identity, err := identity.FlattenSystemAndUserAssignedMap(account.Identity)
		if err != nil {
			return err
		}
		d.Set("identity", identity)

		if props := account.Properties; props != nil {
			if apiProps := props.ApiProperties; apiProps != nil {
				d.Set("qna_runtime_endpoint", apiProps.QnaRuntimeEndpoint)
				d.Set("custom_question_answering_search_service_id", apiProps.QnaAzureSearchEndpointId)
				d.Set("metrics_advisor_aad_client_id", apiProps.AadClientId)
				d.Set("metrics_advisor_aad_tenant_id", apiProps.AadTenantId)
				d.Set("metrics_advisor_super_user_name", apiProps.SuperUser)
				d.Set("metrics_advisor_website_name", apiProps.WebsiteName)
			}
			d.Set("endpoint", props.Endpoint)
			d.Set("custom_subdomain_name", props.CustomSubDomainName)
			if err := d.Set("network_acls", flattenCognitiveAccountNetworkAcls(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_acls` for Cognitive Account %q: %+v", id, err)
			}

			networkInjection, err := flattenCognitiveAccountNetworkInjection(props.NetworkInjections)
			if err != nil {
				return err
			}

			if err := d.Set("network_injection", networkInjection); err != nil {
				return fmt.Errorf("setting `network_injection`: %+v", err)
			}

			d.Set("dynamic_throttling_enabled", pointer.From(props.DynamicThrottlingEnabled))

			d.Set("fqdns", pointer.From(props.AllowedFqdnList))

			// Azure API issue: `AllowProjectManagement` not reset during the rollback, see: https://github.com/Azure/azure-rest-api-specs/issues/38678
			allowProjectManagement := pointer.From(props.AllowProjectManagement)
			if pointer.From(account.Kind) == "OpenAI" {
				allowProjectManagement = false
			}
			d.Set("project_management_enabled", allowProjectManagement)

			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = pointer.From(props.PublicNetworkAccess) == cognitiveservicesaccounts.PublicNetworkAccessEnabled
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)

			if err := d.Set("storage", flattenCognitiveAccountStorage(props.UserOwnedStorage)); err != nil {
				return fmt.Errorf("setting `storages` for Cognitive Account %q: %+v", id, err)
			}
			// lintignore:R001
			d.Set("outbound_network_access_restricted", pointer.From(props.RestrictOutboundNetworkAccess))

			localAuthEnabled := !pointer.From(props.DisableLocalAuth)
			d.Set("local_auth_enabled", localAuthEnabled)

			if localAuthEnabled && fetchAccessKeys {
				keys, err := client.AccountsListKeys(ctx, *id)
				if err != nil {
					return fmt.Errorf("listing the Keys for %s: %+v", *id, err)
				}

				if model := keys.Model; model != nil {
					d.Set("primary_access_key", model.Key1)
					d.Set("secondary_access_key", model.Key2)
				}
			}

			customerManagedKey, err := flattenCognitiveAccountCustomerManagedKey(props.Encryption)
			if err != nil {
				return err
			}

			if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, account.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenCognitiveAccountNetworkAcls(input *cognitiveservicesaccounts.NetworkRuleSet) []any {
	if input == nil {
		return []any{}
	}

	ipRules := make([]any, 0)
	if input.IPRules != nil {
		for _, v := range *input.IPRules {
			ipRules = append(ipRules, v.Value)
		}
	}

	virtualNetworkRules := make([]any, 0)
	if input.VirtualNetworkRules != nil {
		for _, v := range *input.VirtualNetworkRules {
			id := v.Id
			subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
			if err == nil {
				id = subnetId.ID()
			}

			virtualNetworkRules = append(virtualNetworkRules, map[string]any{
				"subnet_id":                            id,
				"ignore_missing_vnet_service_endpoint": *v.IgnoreMissingVnetServiceEndpoint,
			})
		}
	}

	return []any{map[string]any{
		"bypass":                input.Bypass,
		"default_action":        input.DefaultAction,
		"ip_rules":              pluginsdk.NewSet(pluginsdk.HashString, ipRules),
		"virtual_network_rules": virtualNetworkRules,
	}}
}

func flattenCognitiveAccountStorage(input *[]cognitiveservicesaccounts.UserOwnedStorage) []any {
	if input == nil {
		return []any{}
	}
	results := make([]any, 0)
	for _, v := range *input {
		value := make(map[string]any)
		if v.ResourceId != nil {
			value["storage_account_id"] = *v.ResourceId
		}
		if v.IdentityClientId != nil {
			value["identity_client_id"] = *v.IdentityClientId
		}
		results = append(results, value)
	}
	return results
}

func flattenCognitiveAccountCustomerManagedKey(input *cognitiveservicesaccounts.Encryption) ([]any, error) {
	if input == nil || pointer.From(input.KeySource) == cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices {
		return []any{}, nil
	}

	var keyId string
	var identityClientId string
	if props := input.KeyVaultProperties; props != nil {
		keyVaultKeyId, err := keyvault.NewNestedItemID(pointer.From(props.KeyVaultUri), keyvault.NestedItemTypeKey, pointer.From(props.KeyName), pointer.From(props.KeyVersion))
		if err != nil {
			return nil, err
		}
		keyId = keyVaultKeyId.ID()
		if props.IdentityClientId != nil {
			identityClientId = *props.IdentityClientId
		}
	}

	return []any{
		map[string]any{
			"key_vault_key_id":   keyId,
			"identity_client_id": identityClientId,
		},
	}, nil
}

func flattenCognitiveAccountNetworkInjection(input *[]cognitiveservicesaccounts.NetworkInjection) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	results := make([]any, 0)
	for _, v := range *input {
		var subnetId string
		if v.SubnetArmId != nil {
			subnet, err := commonids.ParseSubnetIDInsensitively(*v.SubnetArmId)
			if err != nil {
				return nil, err
			}
			subnetId = subnet.ID()
		}

		results = append(results, map[string]any{
			"scenario":  v.Scenario,
			"subnet_id": subnetId,
		})
	}

	return results, nil
}
