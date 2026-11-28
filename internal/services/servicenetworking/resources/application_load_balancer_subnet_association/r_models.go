// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_subnet_association

type AssociationResourceModel struct {
	Name                      string         `tfschema:"name"`
	ApplicationLoadBalancerId string         `tfschema:"application_load_balancer_id"`
	SubnetId                  string         `tfschema:"subnet_id"`
	Tags                      map[string]any `tfschema:"tags"`
}

func (t Resource) ModelObject() any {
	return &AssociationResourceModel{}
}
