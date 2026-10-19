// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_service

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAccessPolicyEntries(d *pluginsdk.ResourceData) *[]resource.ServiceAccessPolicyEntry {
	accessPolicyObjectIds := d.Get("access_policy_object_ids").(*pluginsdk.Set).List()
	svcAccessPolicyArray := make([]resource.ServiceAccessPolicyEntry, 0)

	for _, objectId := range accessPolicyObjectIds {
		svcAccessPolicyObjectId := resource.ServiceAccessPolicyEntry{ObjectId: objectId.(string)}
		svcAccessPolicyArray = append(svcAccessPolicyArray, svcAccessPolicyObjectId)
	}

	return &svcAccessPolicyArray
}

func expandCorsConfiguration(d *pluginsdk.ResourceData) *resource.ServiceCorsConfigurationInfo {
	corsConfigRaw := d.Get("cors_configuration").([]any)

	if len(corsConfigRaw) == 0 {
		return &resource.ServiceCorsConfigurationInfo{}
	}

	corsConfigAttr := corsConfigRaw[0].(map[string]any)

	allowedOrigins := *pluginsdk.ExpandStringSlice(corsConfigAttr["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *pluginsdk.ExpandStringSlice(corsConfigAttr["allowed_headers"].(*pluginsdk.Set).List())
	allowedMethods := *pluginsdk.ExpandStringSlice(corsConfigAttr["allowed_methods"].([]any))

	return &resource.ServiceCorsConfigurationInfo{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		MaxAge:           pointer.To(int64(corsConfigAttr["max_age_in_seconds"].(int))),
		AllowCredentials: pointer.To(corsConfigAttr["allow_credentials"].(bool)),
	}
}

func expandAuthentication(d *pluginsdk.ResourceData) *resource.ServiceAuthenticationConfigurationInfo {
	authConfigRaw := d.Get("authentication_configuration").([]any)

	if len(authConfigRaw) == 0 {
		return &resource.ServiceAuthenticationConfigurationInfo{}
	}

	authConfigAttr := authConfigRaw[0].(map[string]any)

	return &resource.ServiceAuthenticationConfigurationInfo{
		Authority:         pointer.To(authConfigAttr["authority"].(string)),
		Audience:          pointer.To(authConfigAttr["audience"].(string)),
		SmartProxyEnabled: pointer.To(authConfigAttr["smart_proxy_enabled"].(bool)),
	}
}

func expandsCosmosDBConfiguration(d *pluginsdk.ResourceData) (*resource.ServiceCosmosDbConfigurationInfo, error) {
	cosmosdb := &resource.ServiceCosmosDbConfigurationInfo{
		OfferThroughput: pointer.To(int64(d.Get("cosmosdb_throughput").(int))),
	}

	if keyVaultKeyIDRaw, ok := d.GetOk("cosmosdb_key_vault_key_versionless_id"); ok {
		keyVaultKey, err := keyvault.ParseNestedItemID(keyVaultKeyIDRaw.(string), keyvault.VersionTypeVersionless, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}
		cosmosdb.KeyVaultKeyUri = pointer.To(keyVaultKey.ID())
	}

	return cosmosdb, nil
}
