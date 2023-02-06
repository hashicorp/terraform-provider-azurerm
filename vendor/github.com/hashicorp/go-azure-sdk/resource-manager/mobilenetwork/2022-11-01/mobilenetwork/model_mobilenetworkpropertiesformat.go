package mobilenetwork

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MobileNetworkPropertiesFormat struct {
	ProvisioningState                 *ProvisioningState `json:"provisioningState,omitempty"`
	PublicLandMobileNetworkIdentifier PlmnId             `json:"publicLandMobileNetworkIdentifier"`
	ServiceKey                        *string            `json:"serviceKey,omitempty"`
}
