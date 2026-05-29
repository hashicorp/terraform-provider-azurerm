package proximityplacementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProximityPlacementGroupProperties struct {
	AvailabilitySets            *[]SubResourceWithColocationStatus       `json:"availabilitySets,omitempty"`
	ColocationStatus            *InstanceViewStatus                      `json:"colocationStatus,omitempty"`
	Intent                      *ProximityPlacementGroupPropertiesIntent `json:"intent,omitempty"`
	ProximityPlacementGroupType *ProximityPlacementGroupType             `json:"proximityPlacementGroupType,omitempty"`
	VirtualMachineScaleSets     *[]SubResourceWithColocationStatus       `json:"virtualMachineScaleSets,omitempty"`
	VirtualMachines             *[]SubResourceWithColocationStatus       `json:"virtualMachines,omitempty"`
}
