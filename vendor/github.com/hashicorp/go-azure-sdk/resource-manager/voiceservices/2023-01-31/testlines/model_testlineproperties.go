package testlines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestLineProperties struct {
	PhoneNumber       string             `json:"phoneNumber"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Purpose           TestLinePurpose    `json:"purpose"`
}
