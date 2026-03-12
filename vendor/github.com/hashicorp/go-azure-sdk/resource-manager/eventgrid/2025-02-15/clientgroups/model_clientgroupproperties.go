package clientgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientGroupProperties struct {
	Description       *string                       `json:"description,omitempty"`
	ProvisioningState *ClientGroupProvisioningState `json:"provisioningState,omitempty"`
	Query             *string                       `json:"query,omitempty"`
}
