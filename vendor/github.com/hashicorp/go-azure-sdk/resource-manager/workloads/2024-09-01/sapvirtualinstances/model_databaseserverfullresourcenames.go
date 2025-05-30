package sapvirtualinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseServerFullResourceNames struct {
	AvailabilitySetName *string                        `json:"availabilitySetName,omitempty"`
	LoadBalancer        *LoadBalancerResourceNames     `json:"loadBalancer,omitempty"`
	VirtualMachines     *[]VirtualMachineResourceNames `json:"virtualMachines,omitempty"`
}
