// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack

type LocalRuleStackModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroupName    string `tfschema:"resource_group_name"`
	Location             string `tfschema:"location"`
	AntiSpywareProfile   string `tfschema:"anti_spyware_profile"`
	AntiVirusProfile     string `tfschema:"anti_virus_profile"`
	DNSSubscription      string `tfschema:"dns_subscription"`
	FileBlockingProfile  string `tfschema:"file_blocking_profile"`
	URLFilteringProfile  string `tfschema:"url_filtering_profile"`
	VulnerabilityProfile string `tfschema:"vulnerability_profile"`
	Description          string `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &LocalRuleStackModel{}
}
