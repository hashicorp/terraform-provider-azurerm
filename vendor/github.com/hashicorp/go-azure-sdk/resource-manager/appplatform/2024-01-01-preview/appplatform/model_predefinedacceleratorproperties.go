package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PredefinedAcceleratorProperties struct {
	AcceleratorTags   *[]string                               `json:"acceleratorTags,omitempty"`
	Description       *string                                 `json:"description,omitempty"`
	DisplayName       *string                                 `json:"displayName,omitempty"`
	IconURL           *string                                 `json:"iconUrl,omitempty"`
	ProvisioningState *PredefinedAcceleratorProvisioningState `json:"provisioningState,omitempty"`
	State             *PredefinedAcceleratorState             `json:"state,omitempty"`
}
