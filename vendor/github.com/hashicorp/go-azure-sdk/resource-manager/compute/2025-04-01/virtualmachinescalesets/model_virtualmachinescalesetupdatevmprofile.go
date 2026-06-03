package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetUpdateVMProfile struct {
	BillingProfile           *BillingProfile                             `json:"billingProfile,omitempty"`
	DiagnosticsProfile       *DiagnosticsProfile                         `json:"diagnosticsProfile,omitempty"`
	ExtensionProfile         *VirtualMachineScaleSetExtensionProfile     `json:"extensionProfile,omitempty"`
	HardwareProfile          *VirtualMachineScaleSetHardwareProfile      `json:"hardwareProfile,omitempty"`
	LicenseType              *string                                     `json:"licenseType,omitempty"`
	NetworkProfile           *VirtualMachineScaleSetUpdateNetworkProfile `json:"networkProfile,omitempty"`
	OsProfile                *VirtualMachineScaleSetUpdateOSProfile      `json:"osProfile,omitempty"`
	ScheduledEventsProfile   *ScheduledEventsProfile                     `json:"scheduledEventsProfile,omitempty"`
	SecurityPostureReference *SecurityPostureReferenceUpdate             `json:"securityPostureReference,omitempty"`
	SecurityProfile          *SecurityProfile                            `json:"securityProfile,omitempty"`
	StorageProfile           *VirtualMachineScaleSetUpdateStorageProfile `json:"storageProfile,omitempty"`
	UserData                 *string                                     `json:"userData,omitempty"`
}
