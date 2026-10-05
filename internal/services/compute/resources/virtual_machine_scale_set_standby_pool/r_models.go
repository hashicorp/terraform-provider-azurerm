// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_standby_pool

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbyvirtualmachinepools"
)

type VirtualMachineScaleSetStandbyPoolModel struct {
	Name                             string                                                    `tfschema:"name"`
	ResourceGroupName                string                                                    `tfschema:"resource_group_name"`
	Location                         string                                                    `tfschema:"location"`
	AttachedVirtualMachineScaleSetId string                                                    `tfschema:"attached_virtual_machine_scale_set_id"`
	ElasticityProfile                []VirtualMachineScaleSetStandbyPoolElasticityProfileModel `tfschema:"elasticity_profile"`
	VirtualMachineState              standbyvirtualmachinepools.VirtualMachineState            `tfschema:"virtual_machine_state"`
	Tags                             map[string]string                                         `tfschema:"tags"`
}

type VirtualMachineScaleSetStandbyPoolElasticityProfileModel struct {
	MaxReadyCapacity int64 `tfschema:"max_ready_capacity"`
	MinReadyCapacity int64 `tfschema:"min_ready_capacity"`
}

func (r Resource) ModelObject() any {
	return &VirtualMachineScaleSetStandbyPoolModel{}
}
