// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type MaintenanceAssignmentVirtualMachineScaleSetId struct {
	VirtualMachineScaleSetId    *commonids.VirtualMachineScaleSetId
	VirtualMachineScaleSetIdRaw string
	Name                        string
}
