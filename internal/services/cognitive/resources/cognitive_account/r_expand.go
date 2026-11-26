// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandCognitiveAccountNetworkAcls(d *pluginsdk.ResourceData) (*cognitiveservicesaccounts.NetworkRuleSet, []string) {
	input := d.Get("network_acls").([]any)
	subnetIds := make([]string, 0)
	if len(input) == 0 || input[0] == nil {
		return nil, subnetIds
	}

	v := input[0].(map[string]any)

	ipRulesRaw := v["ip_rules"].(*pluginsdk.Set)
	ipRules := make([]cognitiveservicesaccounts.IPRule, 0)

	for _, v := range ipRulesRaw.List() {
		rule := cognitiveservicesaccounts.IPRule{
			Value: v.(string),
		}
		ipRules = append(ipRules, rule)
	}

	networkRules := make([]cognitiveservicesaccounts.VirtualNetworkRule, 0)
	networkRulesRaw := v["virtual_network_rules"]
	for _, v := range networkRulesRaw.(*pluginsdk.Set).List() {
		value := v.(map[string]any)
		subnetId := value["subnet_id"].(string)
		subnetIds = append(subnetIds, subnetId)
		rule := cognitiveservicesaccounts.VirtualNetworkRule{
			Id:                               subnetId,
			IgnoreMissingVnetServiceEndpoint: pointer.To(value["ignore_missing_vnet_service_endpoint"].(bool)),
		}
		networkRules = append(networkRules, rule)
	}

	ruleSet := cognitiveservicesaccounts.NetworkRuleSet{
		DefaultAction:       pointer.ToEnum[cognitiveservicesaccounts.NetworkRuleAction](v["default_action"].(string)),
		IPRules:             &ipRules,
		VirtualNetworkRules: &networkRules,
	}

	if b, ok := d.GetOk("network_acls.0.bypass"); ok && b != "" {
		ruleSet.Bypass = pointer.ToEnum[cognitiveservicesaccounts.ByPassSelection](v["bypass"].(string))
	}

	return &ruleSet, subnetIds
}

func expandCognitiveAccountStorage(input []any) *[]cognitiveservicesaccounts.UserOwnedStorage {
	if len(input) == 0 {
		return nil
	}
	results := make([]cognitiveservicesaccounts.UserOwnedStorage, 0)
	for _, v := range input {
		value := v.(map[string]any)
		results = append(results, cognitiveservicesaccounts.UserOwnedStorage{
			ResourceId:       pointer.To(value["storage_account_id"].(string)),
			IdentityClientId: pointer.To(value["identity_client_id"].(string)),
		})
	}
	return &results
}

func expandCognitiveAccountAPIProperties(d *pluginsdk.ResourceData) *cognitiveservicesaccounts.ApiProperties {
	props := cognitiveservicesaccounts.ApiProperties{}
	kind := d.Get("kind")

	if kind == "QnAMaker" {
		if v, ok := d.GetOk("qna_runtime_endpoint"); ok {
			props.QnaRuntimeEndpoint = pointer.To(v.(string))
		}
	}

	if kind == "TextAnalytics" {
		if v, ok := d.GetOk("custom_question_answering_search_service_id"); ok {
			props.QnaAzureSearchEndpointId = pointer.To(v.(string))
		}
		if v, ok := d.GetOk("custom_question_answering_search_service_key"); ok {
			props.QnaAzureSearchEndpointKey = pointer.To(v.(string))
		}
	}

	if kind == "MetricsAdvisor" {
		if v, ok := d.GetOk("metrics_advisor_aad_client_id"); ok {
			props.AadClientId = pointer.To(v.(string))
		}
		if v, ok := d.GetOk("metrics_advisor_aad_tenant_id"); ok {
			props.AadTenantId = pointer.To(v.(string))
		}
		if v, ok := d.GetOk("metrics_advisor_super_user_name"); ok {
			props.SuperUser = pointer.To(v.(string))
		}
		if v, ok := d.GetOk("metrics_advisor_website_name"); ok {
			props.WebsiteName = pointer.To(v.(string))
		}
	}

	return &props
}

func expandCognitiveAccountCustomerManagedKey(input []any) *cognitiveservicesaccounts.Encryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	keyId, _ := keyvault.ParseNestedItemID(v["key_vault_key_id"].(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if keyId == nil {
		return nil
	}

	var identity string
	if value := v["identity_client_id"]; value != nil && value != "" {
		identity = value.(string)
	}

	return &cognitiveservicesaccounts.Encryption{
		KeySource: pointer.To(cognitiveservicesaccounts.KeySourceMicrosoftPointKeyVault),
		KeyVaultProperties: &cognitiveservicesaccounts.KeyVaultProperties{
			KeyName:          pointer.To(keyId.Name),
			KeyVersion:       pointer.To(keyId.Version),
			KeyVaultUri:      pointer.To(keyId.KeyVaultBaseURL),
			IdentityClientId: pointer.To(identity),
		},
	}
}

func expandCognitiveAccountNetworkInjection(input []any) *[]cognitiveservicesaccounts.NetworkInjection {
	if len(input) == 0 {
		return nil
	}

	results := make([]cognitiveservicesaccounts.NetworkInjection, 0)
	for _, v := range input {
		m := v.(map[string]any)

		var subnetId *string
		if m["subnet_id"] != nil && m["subnet_id"] != "" {
			subnetId = pointer.To(m["subnet_id"].(string))
		}

		results = append(results, cognitiveservicesaccounts.NetworkInjection{
			Scenario:    pointer.ToEnum[cognitiveservicesaccounts.ScenarioType](m["scenario"].(string)),
			SubnetArmId: subnetId,
		})
	}

	return &results
}
