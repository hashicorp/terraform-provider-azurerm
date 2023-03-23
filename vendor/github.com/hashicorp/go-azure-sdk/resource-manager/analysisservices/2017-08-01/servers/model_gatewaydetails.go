package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayDetails struct {
	DmtsClusterUri    *string `json:"dmtsClusterUri,omitempty"`
	GatewayObjectId   *string `json:"gatewayObjectId,omitempty"`
	GatewayResourceId *string `json:"gatewayResourceId,omitempty"`
}
