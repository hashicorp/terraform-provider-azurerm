package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualApplianceNicProperties struct {
	InstanceName     *string            `json:"instanceName,omitempty"`
	Name             *string            `json:"name,omitempty"`
	NicType          *NicTypeInResponse `json:"nicType,omitempty"`
	PrivateIPAddress *string            `json:"privateIpAddress,omitempty"`
	PublicIPAddress  *string            `json:"publicIpAddress,omitempty"`
}
