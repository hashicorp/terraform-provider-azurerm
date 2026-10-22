// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster_firewall_rule

type MongoClusterFirewallRuleResourceModel struct {
	MongoClusterId string `tfschema:"mongo_cluster_id"`
	Name           string `tfschema:"name"`
	EndIpAddress   string `tfschema:"end_ip_address"`
	StartIpAddress string `tfschema:"start_ip_address"`
}

func (r Resource) ModelObject() any {
	return &MongoClusterFirewallRuleResourceModel{}
}
