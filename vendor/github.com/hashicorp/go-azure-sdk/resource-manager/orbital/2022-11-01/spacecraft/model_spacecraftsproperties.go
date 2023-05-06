package spacecraft

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpacecraftsProperties struct {
	Links             []SpacecraftLink   `json:"links"`
	NoradId           *string            `json:"noradId,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	TitleLine         string             `json:"titleLine"`
	TleLine1          string             `json:"tleLine1"`
	TleLine2          string             `json:"tleLine2"`
}
