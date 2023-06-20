package lab

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineProfile struct {
	AdditionalCapabilities *VirtualMachineAdditionalCapabilities `json:"additionalCapabilities,omitempty"`
	AdminUser              Credentials                           `json:"adminUser"`
	CreateOption           CreateOption                          `json:"createOption"`
	ImageReference         ImageReference                        `json:"imageReference"`
	NonAdminUser           *Credentials                          `json:"nonAdminUser,omitempty"`
	OsType                 *OsType                               `json:"osType,omitempty"`
	Sku                    Sku                                   `json:"sku"`
	UsageQuota             string                                `json:"usageQuota"`
	UseSharedPassword      *EnableState                          `json:"useSharedPassword,omitempty"`
}
