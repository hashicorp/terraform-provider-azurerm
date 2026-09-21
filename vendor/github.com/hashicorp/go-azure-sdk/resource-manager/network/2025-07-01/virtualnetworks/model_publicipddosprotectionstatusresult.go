package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicIPDdosProtectionStatusResult struct {
	DdosProtectionPlanId *string              `json:"ddosProtectionPlanId,omitempty"`
	IsWorkloadProtected  *IsWorkloadProtected `json:"isWorkloadProtected,omitempty"`
	PublicIPAddress      *string              `json:"publicIpAddress,omitempty"`
	PublicIPAddressId    *string              `json:"publicIpAddressId,omitempty"`
}
