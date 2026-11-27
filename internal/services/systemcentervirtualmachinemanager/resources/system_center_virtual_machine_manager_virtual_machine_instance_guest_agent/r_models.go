// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance_guest_agent

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel struct {
	ScopedResourceId   string `tfschema:"scoped_resource_id"`
	Username           string `tfschema:"username"`
	Password           string `tfschema:"password"`
	ProvisioningAction string `tfschema:"provisioning_action"`
}

func (r Resource) ModelObject() any {
	return &SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentModel{}
}
