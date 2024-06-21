package availabilitysets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilitySetProperties struct {
	AvailabilitySetName *string                    `json:"availabilitySetName,omitempty"`
	ProvisioningState   *ResourceProvisioningState `json:"provisioningState,omitempty"`
	VMmServerId         *string                    `json:"vmmServerId,omitempty"`
}
