// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var ValidateWebApplicationFirewallConfigurationRuleGroupName = validation.StringInSlice([]string{
	"BadBots",
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
	"GoodBots",
	"Known-CVEs",
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
	"UnknownBots",
}, false)

var ValidateWebApplicationFirewallConfigurationRuleSetVersion = validation.StringInSlice([]string{
	"0.1",
	"1.0",
	"2.2.9",
	"3.0",
	"3.1",
	"3.2",
}, false)

var ValidateWebApplicationFirewallConfigurationRuleSetType = validation.StringInSlice([]string{
	"OWASP",
	"Microsoft_BotManagerRuleSet",
}, false)

// https://learn.microsoft.com/en-us/azure/web-application-firewall/ag/application-gateway-crs-rulegroups-rules?tabs=drs21
var ValidateWebApplicationFirewallPolicyRuleGroupName = validation.StringInSlice([]string{
	"APPLICATION-ATTACK-LFI",
	"APPLICATION-ATTACK-NodeJS",
	"APPLICATION-ATTACK-PHP",
	"APPLICATION-ATTACK-RCE",
	"APPLICATION-ATTACK-RFI",
	"APPLICATION-ATTACK-SQLI",
	"APPLICATION-ATTACK-SESSION-FIXATION",
	"APPLICATION-ATTACK-SESSION-JAVA",
	"APPLICATION-ATTACK-XSS",
	"BadBots",
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
	"GoodBots",
	"Known-CVEs",
	"METHOD-ENFORCEMENT",
	"MS-ThreatIntel-AppSec",
	"MS-ThreatIntel-CVEs",
	"MS-ThreatIntel-SQLI",
	"MS-ThreatIntel-WebShells",
	"PROTOCOL-ATTACK",
	"PROTOCOL-ENFORCEMENT",
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
	"UnknownBots",
}, false)

var ValidateWebApplicationFirewallPolicyRuleSetVersion = validation.StringInSlice([]string{
	"0.1",
	"1.0",
	"2.1",
	"2.2.9",
	"3.0",
	"3.1",
	"3.2",
}, false)

var ValidateWebApplicationFirewallPolicyRuleSetType = validation.StringInSlice([]string{
	"OWASP",
	"Microsoft_BotManagerRuleSet",
	"Microsoft_DefaultRuleSet",
}, false)

var ValidateWebApplicationFirewallPolicyExclusionRuleSetVersion = validation.StringInSlice([]string{
	"3.2",
	"2.1",
}, false)

var ValidateWebApplicationFirewallPolicyExclusionRuleSetType = validation.StringInSlice([]string{
	"OWASP",
	"Microsoft_DefaultRuleSet",
}, false)
