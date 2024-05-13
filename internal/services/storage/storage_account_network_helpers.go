package storage

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
)

// TODO: remove the `NetworkRule` name in time

func expandStorageAccountNetworkRuleBypass(input []interface{}) *storageaccounts.Bypass {
	if len(input) == 0 {
		return nil
	}

	output := make([]string, 0)
	for _, item := range input {
		output = append(output, item.(string))
	}
	return pointer.To(storageaccounts.Bypass(strings.Join(output, ", ")))
}

func flattenStorageAccountNetworkRulesBypass(input *storageaccounts.Bypass) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		// Whilst this is an Enum it's actually a CSV containing the enum but it's exposed as a regular string
		// as such we need to manually re-case this
		components := strings.Split(string(*input), ", ")

		vals := map[string]string{}
		for _, v := range storageaccounts.PossibleValuesForBypass() {
			vals[strings.ToLower(v)] = v
		}
		for _, item := range components {
			val, ok := vals[strings.ToLower(item)]
			if !ok {
				// otherwise append the direct casing if it's an unknown value
				val = item
			}

			output = append(output, val)
		}
	}

	return output
}

func expandStorageAccountNetworkRuleIPRules(input []interface{}) *[]storageaccounts.IPRule {
	output := make([]storageaccounts.IPRule, 0)

	for _, item := range input {
		output = append(output, storageaccounts.IPRule{
			Action: pointer.To(storageaccounts.ActionAllow),
			Value:  item.(string),
		})
	}

	return &output
}

func flattenStorageAccountNetworkRulesIPRules(input *[]storageaccounts.IPRule) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, item.Value)
		}
	}

	return output
}

func expandStorageAccountNetworkRuleVirtualRules(input []interface{}) *[]storageaccounts.VirtualNetworkRule {
	output := make([]storageaccounts.VirtualNetworkRule, 0)

	for _, item := range input {
		output = append(output, storageaccounts.VirtualNetworkRule{
			Id:     item.(string),
			Action: pointer.To(storageaccounts.ActionAllow),
		})
	}

	return &output
}

func flattenStorageAccountNetworkRulesVirtualNetworkRules(input *[]storageaccounts.VirtualNetworkRule) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			output = append(output, item.Id)
		}
	}

	return output
}

func expandStorageAccountNetworkRulesPrivateLinkAccess(input []interface{}, tenantId string) *[]storageaccounts.ResourceAccessRule {
	output := make([]storageaccounts.ResourceAccessRule, 0)
	for _, raw := range input {
		item := raw.(map[string]interface{})
		rule := storageaccounts.ResourceAccessRule{
			ResourceId: pointer.To(item["endpoint_resource_id"].(string)),
			TenantId:   pointer.To(tenantId),
		}
		if v := item["endpoint_tenant_id"].(string); v != "" {
			rule.TenantId = pointer.To(v)
		}
		output = append(output, rule)
	}

	return &output
}

func flattenStorageAccountNetworkRulesPrivateLinkAccess(input *[]storageaccounts.ResourceAccessRule) []interface{} {
	output := make([]interface{}, 0)

	for _, item := range *input {
		output = append(output, map[string]interface{}{
			"endpoint_resource_id": pointer.From(item.ResourceId),
			"endpoint_tenant_id":   pointer.From(item.TenantId),
		})
	}

	return output
}
