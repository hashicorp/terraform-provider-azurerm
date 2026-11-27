package environments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentProperties struct {
	CustomProperties *interface{}       `json:"customProperties,omitempty"`
	Description      *string            `json:"description,omitempty"`
	Kind             EnvironmentKind    `json:"kind"`
	Onboarding       *Onboarding        `json:"onboarding,omitempty"`
	Server           *EnvironmentServer `json:"server,omitempty"`
	Title            string             `json:"title"`
}
