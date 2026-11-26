package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolImage struct {
	Aliases            *[]string      `json:"aliases,omitempty"`
	Buffer             *string        `json:"buffer,omitempty"`
	EphemeralType      *EphemeralType `json:"ephemeralType,omitempty"`
	IsEphemeral        *bool          `json:"isEphemeral,omitempty"`
	ResourceId         *string        `json:"resourceId,omitempty"`
	WellKnownImageName *string        `json:"wellKnownImageName,omitempty"`
}
