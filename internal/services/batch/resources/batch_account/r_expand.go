// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
)

func expandEncryption(e []interface{}) *batchaccount.EncryptionProperties {
	defaultEnc := batchaccount.EncryptionProperties{
		KeySource: pointer.To(batchaccount.KeySourceMicrosoftPointBatch),
	}

	if len(e) == 0 || e[0] == nil {
		return &defaultEnc
	}

	v := e[0].(map[string]interface{})
	encryptionProperty := batchaccount.EncryptionProperties{
		KeySource: pointer.To(batchaccount.KeySourceMicrosoftPointKeyVault),
		KeyVaultProperties: &batchaccount.KeyVaultProperties{
			KeyIdentifier: pointer.To(v["key_vault_key_id"].(string)),
		},
	}

	return &encryptionProperty
}

func expandAllowedAuthenticationModes(input []interface{}) *[]batchaccount.AuthenticationMode {
	if len(input) == 0 {
		return nil
	}

	allowedAuthModes := make([]batchaccount.AuthenticationMode, 0)
	for _, mode := range input {
		allowedAuthModes = append(allowedAuthModes, batchaccount.AuthenticationMode(mode.(string)))
	}
	return &allowedAuthModes
}

func expandBatchAccountNetworkProfile(input []interface{}) *batchaccount.NetworkProfile {
	if len(input) == 0 || input[0] == nil {
		return &batchaccount.NetworkProfile{}
	}

	networkProfile := input[0].(map[string]interface{})
	return &batchaccount.NetworkProfile{
		AccountAccess:        expandBatchAccountEndpointAccessProfile(networkProfile["account_access"].([]interface{})),
		NodeManagementAccess: expandBatchAccountEndpointAccessProfile(networkProfile["node_management_access"].([]interface{})),
	}
}

func expandBatchAccountEndpointAccessProfile(input []interface{}) *batchaccount.EndpointAccessProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	accessProfile := input[0].(map[string]interface{})

	ipRulesRaw := accessProfile["ip_rule"].([]interface{})
	ipRules := make([]batchaccount.IPRule, 0)
	for _, ipRule := range ipRulesRaw {
		ipRuleRaw := ipRule.(map[string]interface{})
		ipRules = append(ipRules, batchaccount.IPRule{
			Action: batchaccount.IPRuleAction(ipRuleRaw["action"].(string)),
			Value:  ipRuleRaw["ip_range"].(string),
		})
	}

	return &batchaccount.EndpointAccessProfile{
		DefaultAction: batchaccount.EndpointAccessDefaultAction(accessProfile["default_action"].(string)),
		IPRules:       pointer.To(ipRules),
	}
}
