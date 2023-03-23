package spacecraft

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpacecraftsProperties struct {
	Links             *[]SpacecraftLink  `json:"links,omitempty"`
	NoradId           string             `json:"noradId"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	TitleLine         *string            `json:"titleLine,omitempty"`
	TleLine1          *string            `json:"tleLine1,omitempty"`
	TleLine2          *string            `json:"tleLine2,omitempty"`
}
