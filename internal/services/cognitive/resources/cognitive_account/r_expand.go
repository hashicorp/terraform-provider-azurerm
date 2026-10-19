// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"fmt"

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

func expandCognitiveAccountAPIProperties(d *pluginsdk.ResourceData) (*cognitiveservicesaccounts.ApiProperties, error) {
	props := cognitiveservicesaccounts.ApiProperties{}
	kind := d.Get("kind")
	if kind == "QnAMaker" {
		if v, ok := d.GetOk("qna_runtime_endpoint"); ok && v != "" {
			props.QnaRuntimeEndpoint = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`")
		}
	}
	if v, ok := d.GetOk("custom_question_answering_search_service_id"); ok {
		if kind == "TextAnalytics" {
			props.QnaAzureSearchEndpointId = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("the Search Service ID `custom_question_answering_search_service_id` can only be set when kind is set to `TextAnalytics`")
		}
	}
	if v, ok := d.GetOk("custom_question_answering_search_service_key"); ok {
		if kind == "TextAnalytics" {
			props.QnaAzureSearchEndpointKey = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("the Search Service Key `custom_question_answering_search_service_key` can only be set when kind is set to `TextAnalytics`")
		}
	}

	if v, ok := d.GetOk("metrics_advisor_aad_client_id"); ok {
		if kind == "MetricsAdvisor" {
			props.AadClientId = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_aad_client_id can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_aad_tenant_id"); ok {
		if kind == "MetricsAdvisor" {
			props.AadTenantId = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_aad_tenant_id can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_super_user_name"); ok {
		if kind == "MetricsAdvisor" {
			props.SuperUser = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_super_user_name can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_website_name"); ok {
		if kind == "MetricsAdvisor" {
			props.WebsiteName = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_website_name can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	return &props, nil
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
