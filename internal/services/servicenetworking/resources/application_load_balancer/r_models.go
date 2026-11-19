// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer

type ApplicationLoadBalancerModel struct {
	Name                          string            `tfschema:"name"`
	ResourceGroupName             string            `tfschema:"resource_group_name"`
	Location                      string            `tfschema:"location"`
	PrimaryConfigurationEndpoints string            `tfschema:"primary_configuration_endpoint"`
	Tags                          map[string]string `tfschema:"tags"`
}

func (t Resource) ModelObject() any {
	return &ApplicationLoadBalancerModel{}
}
