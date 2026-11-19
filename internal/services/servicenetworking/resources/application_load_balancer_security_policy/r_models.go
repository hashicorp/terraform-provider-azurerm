// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_security_policy

type SecurityPoliciesModel struct {
	ApplicationLoadBalancerId      string            `tfschema:"application_load_balancer_id"`
	Location                       string            `tfschema:"location"`
	Name                           string            `tfschema:"name"`
	WebApplicationFirewallPolicyId string            `tfschema:"web_application_firewall_policy_id"`
	Tags                           map[string]string `tfschema:"tags"`
}

func (f Resource) ModelObject() any {
	return &SecurityPoliciesModel{}
}
