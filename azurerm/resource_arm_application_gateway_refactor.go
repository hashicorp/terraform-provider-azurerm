package azurerm

import (
	"bytes"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func expandApplicationGatewayWafConfig(d *schema.ResourceData) *network.ApplicationGatewayWebApplicationFirewallConfiguration {
	wafSet := d.Get("waf_configuration").(*schema.Set).List()
	waf := wafSet[0].(map[string]interface{})

	enabled := waf["enabled"].(bool)
	mode := waf["firewall_mode"].(string)
	rulesettype := waf["rule_set_type"].(string)
	rulesetversion := waf["rule_set_version"].(string)

	return &network.ApplicationGatewayWebApplicationFirewallConfiguration{
		Enabled:        &enabled,
		FirewallMode:   network.ApplicationGatewayFirewallMode(mode),
		RuleSetType:    &rulesettype,
		RuleSetVersion: &rulesetversion,
	}
}

func flattenApplicationGatewayWafConfig(waf *network.ApplicationGatewayWebApplicationFirewallConfiguration) []interface{} {
	result := make(map[string]interface{})

	result["enabled"] = *waf.Enabled
	result["firewall_mode"] = string(waf.FirewallMode)
	result["rule_set_type"] = waf.RuleSetType
	result["rule_set_version"] = waf.RuleSetVersion

	return []interface{}{result}
}

// TODO: can this be removed?
func hashApplicationGatewayWafConfig(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%t-", m["enabled"].(bool)))
	buf.WriteString(fmt.Sprintf("%s-", m["firewall_mode"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", *m["rule_set_type"].(*string)))
	buf.WriteString(fmt.Sprintf("%s-", *m["rule_set_version"].(*string)))

	return hashcode.String(buf.String())
}
