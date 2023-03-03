package capacityreservationgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityReservationGroupProperties struct {
	CapacityReservations      *[]SubResourceReadOnly                `json:"capacityReservations,omitempty"`
	InstanceView              *CapacityReservationGroupInstanceView `json:"instanceView,omitempty"`
	VirtualMachinesAssociated *[]SubResourceReadOnly                `json:"virtualMachinesAssociated,omitempty"`
}
