package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedKeyProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	SharedKey         *string            `json:"sharedKey,omitempty"`
	SharedKeyLength   *int64             `json:"sharedKeyLength,omitempty"`
}
