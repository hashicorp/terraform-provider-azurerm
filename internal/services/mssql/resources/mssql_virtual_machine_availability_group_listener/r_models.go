// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_availability_group_listener

type MsSqlVirtualMachineAvailabilityGroupListenerModel struct {
	Name                     string `tfschema:"name"`
	SqlVirtualMachineGroupId string `tfschema:"sql_virtual_machine_group_id"`
	AvailabilityGroupName    string `tfschema:"availability_group_name"`

	Port                       int64                                                                    `tfschema:"port"`
	LoadBalancerConfiguration  []LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener  `tfschema:"load_balancer_configuration"`
	MultiSubnetIpConfiguration []MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener `tfschema:"multi_subnet_ip_configuration"`
	Replica                    []ReplicaMsSqlVirtualMachineAvailabilityGroupListener                    `tfschema:"replica"`
}

func (r Resource) ModelObject() any {
	return &MsSqlVirtualMachineAvailabilityGroupListenerModel{}
}
