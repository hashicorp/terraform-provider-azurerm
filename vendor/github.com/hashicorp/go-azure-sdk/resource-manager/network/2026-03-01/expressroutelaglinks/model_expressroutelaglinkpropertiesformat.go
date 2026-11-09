package expressroutelaglinks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagLinkPropertiesFormat struct {
	AdminState        *ExpressRouteLinkAdminState   `json:"adminState,omitempty"`
	InterfaceName     *string                       `json:"interfaceName,omitempty"`
	MacSecConfig      *ExpressRouteLinkMacSecConfig `json:"macSecConfig,omitempty"`
	Members           *[]ExpressRouteLagMember      `json:"members,omitempty"`
	ProvisioningState *ProvisioningState            `json:"provisioningState,omitempty"`
	RouterName        *string                       `json:"routerName,omitempty"`
}
