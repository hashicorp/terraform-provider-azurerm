package expressrouteports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkPropertiesFormat struct {
	AdminState        *ExpressRouteLinkAdminState    `json:"adminState,omitempty"`
	ColoLocation      *string                        `json:"coloLocation,omitempty"`
	ConnectorType     *ExpressRouteLinkConnectorType `json:"connectorType,omitempty"`
	InterfaceName     *string                        `json:"interfaceName,omitempty"`
	MacSecConfig      *ExpressRouteLinkMacSecConfig  `json:"macSecConfig,omitempty"`
	PatchPanelId      *string                        `json:"patchPanelId,omitempty"`
	ProvisioningState *ProvisioningState             `json:"provisioningState,omitempty"`
	RackId            *string                        `json:"rackId,omitempty"`
	RouterName        *string                        `json:"routerName,omitempty"`
}
