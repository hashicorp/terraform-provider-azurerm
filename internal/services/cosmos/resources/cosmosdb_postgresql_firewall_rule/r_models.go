// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_firewall_rule

type CosmosDbPostgreSQLFirewallRuleResourceModel struct {
	Name           string `tfschema:"name"`
	ClusterId      string `tfschema:"cluster_id"`
	EndIPAddress   string `tfschema:"end_ip_address"`
	StartIPAddress string `tfschema:"start_ip_address"`
}

func (r Resource) ModelObject() any {
	return &CosmosDbPostgreSQLFirewallRuleResourceModel{}
}
