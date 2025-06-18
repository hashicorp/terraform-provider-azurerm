package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubstituteFromPatchDefinition struct {
	Kind     *string `json:"kind,omitempty"`
	Name     *string `json:"name,omitempty"`
	Optional *bool   `json:"optional,omitempty"`
}
