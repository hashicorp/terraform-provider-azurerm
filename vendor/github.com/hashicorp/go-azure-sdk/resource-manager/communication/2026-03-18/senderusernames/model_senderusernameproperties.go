package senderusernames

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SenderUsernameProperties struct {
	DataLocation      *string            `json:"dataLocation,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Username          string             `json:"username"`
}
