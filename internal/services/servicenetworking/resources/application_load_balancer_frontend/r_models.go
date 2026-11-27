// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_frontend

type FrontendsModel struct {
	Name                      string         `tfschema:"name"`
	ApplicationLoadBalancerId string         `tfschema:"application_load_balancer_id"`
	Fqdn                      string         `tfschema:"fully_qualified_domain_name"`
	Tags                      map[string]any `tfschema:"tags"`
}

func (f Resource) ModelObject() any {
	return &FrontendsModel{}
}
