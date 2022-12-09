package slice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SlicePropertiesFormat struct {
	Description       *string            `json:"description,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Snssai            Snssai             `json:"snssai"`
}
