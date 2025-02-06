package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachinePublicIPAddressDnsSettingsConfiguration struct {
	DomainNameLabel      string                     `json:"domainNameLabel"`
	DomainNameLabelScope *DomainNameLabelScopeTypes `json:"domainNameLabelScope,omitempty"`
}
