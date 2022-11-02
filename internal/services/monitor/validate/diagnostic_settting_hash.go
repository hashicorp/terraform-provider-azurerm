package validate

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ResourceMonitorDiagnosticLogSettingHash(input interface{}) int {
	var buf bytes.Buffer
	if rawData, ok := input.(map[string]interface{}); ok {
		if category, ok := rawData["category"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", category.(string)))
		}
		if categoryGroup, ok := rawData["category_group"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", categoryGroup.(string)))
		}
		if enabled, ok := rawData["enabled"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", enabled.(bool)))
		}
		if policy, ok := rawData["retention_policy"].(map[string]interface{}); ok {
			if policyEnabled, ok := policy["enabled"]; ok {
				buf.WriteString(fmt.Sprintf("%t-", policyEnabled.(bool)))
			}
			if days, ok := policy["days"]; ok {
				buf.WriteString(fmt.Sprintf("%d-", days.(int)))
			}
		}
	}
	return pluginsdk.HashString(buf.String())
}

func ResourceMonitorDiagnosticMetricsSettingHash(input interface{}) int {
	var buf bytes.Buffer
	if rawData, ok := input.(map[string]interface{}); ok {
		if category, ok := rawData["category"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", category.(string)))
		}
		if enabled, ok := rawData["enabled"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", enabled.(bool)))
		}
		if policy, ok := rawData["retention_policy"].(map[string]interface{}); ok {
			if policyEnabled, ok := policy["enabled"]; ok {
				buf.WriteString(fmt.Sprintf("%t-", policyEnabled.(bool)))
			}
			if days, ok := policy["days"]; ok {
				buf.WriteString(fmt.Sprintf("%d-", days.(int)))
			}
		}
	}
	return pluginsdk.HashString(buf.String())
}
