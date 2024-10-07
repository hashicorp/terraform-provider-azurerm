package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstancePropertiesSecurityProfile struct {
	EnableTPM    *bool                                                        `json:"enableTPM,omitempty"`
	SecurityType *SecurityTypes                                               `json:"securityType,omitempty"`
	UefiSettings *VirtualMachineInstancePropertiesSecurityProfileUefiSettings `json:"uefiSettings,omitempty"`
}
