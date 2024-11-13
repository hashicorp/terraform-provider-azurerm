package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityProfileProperties struct {
	TargetVMConfidentialEncryption *SecurityConfiguration `json:"targetVmConfidentialEncryption,omitempty"`
	TargetVMMonitoring             *SecurityConfiguration `json:"targetVmMonitoring,omitempty"`
	TargetVMSecureBoot             *SecurityConfiguration `json:"targetVmSecureBoot,omitempty"`
	TargetVMSecurityType           *SecurityType          `json:"targetVmSecurityType,omitempty"`
	TargetVMTpm                    *SecurityConfiguration `json:"targetVmTpm,omitempty"`
}
