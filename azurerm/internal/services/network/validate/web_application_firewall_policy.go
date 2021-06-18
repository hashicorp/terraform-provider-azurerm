package validate

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

var ValidateWebApplicationFirewallPolicyRuleGroupName = validation.StringInSlice([]string{
	"crs_20_protocol_violations",
	"crs_21_protocol_anomalies",
	"crs_23_request_limits",
	"crs_30_http_policy",
	"crs_35_bad_robots",
	"crs_40_generic_attacks",
	"crs_41_sql_injection_attacks",
	"crs_41_xss_attacks",
	"crs_42_tight_security",
	"crs_45_trojans",
	"General",
	"REQUEST-911-METHOD-ENFORCEMENT",
	"REQUEST-913-SCANNER-DETECTION",
	"REQUEST-920-PROTOCOL-ENFORCEMENT",
	"REQUEST-921-PROTOCOL-ATTACK",
	"REQUEST-930-APPLICATION-ATTACK-LFI",
	"REQUEST-931-APPLICATION-ATTACK-RFI",
	"REQUEST-932-APPLICATION-ATTACK-RCE",
	"REQUEST-933-APPLICATION-ATTACK-PHP",
	"REQUEST-941-APPLICATION-ATTACK-XSS",
	"REQUEST-942-APPLICATION-ATTACK-SQLI",
	"REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION",
	"REQUEST-944-APPLICATION-ATTACK-JAVA",
}, false)

var ValidateWebApplicationFirewallPolicyRuleSetVersion = validation.StringInSlice([]string{
	"0.1",
	"1.0",
	"2.2.9",
	"3.0",
	"3.1",
	"3.2",
}, false)

var ValidateWebApplicationFirewallPolicyRuleSetType = validation.StringInSlice([]string{
	"OWASP",
	"Microsoft_BotManagerRuleSet",
}, false)
