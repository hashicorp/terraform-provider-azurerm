package dedicatedhostgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostGroupProperties struct {
	AdditionalCapabilities    *DedicatedHostGroupPropertiesAdditionalCapabilities `json:"additionalCapabilities,omitempty"`
	Hosts                     *[]SubResourceReadOnly                              `json:"hosts,omitempty"`
	InstanceView              *DedicatedHostGroupInstanceView                     `json:"instanceView,omitempty"`
	PlatformFaultDomainCount  int64                                               `json:"platformFaultDomainCount"`
	SupportAutomaticPlacement *bool                                               `json:"supportAutomaticPlacement,omitempty"`
}
