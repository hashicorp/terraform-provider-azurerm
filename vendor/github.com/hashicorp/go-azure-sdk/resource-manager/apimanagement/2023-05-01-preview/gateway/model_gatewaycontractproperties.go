package gateway

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayContractProperties struct {
	Description  *string                       `json:"description,omitempty"`
	LocationData *ResourceLocationDataContract `json:"locationData,omitempty"`
}
