package virtualmachineruncommands

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunCommandDocument struct {
	Description string                           `json:"description"`
	Id          string                           `json:"id"`
	Label       string                           `json:"label"`
	OsType      OperatingSystemTypes             `json:"osType"`
	Parameters  *[]RunCommandParameterDefinition `json:"parameters,omitempty"`
	Schema      string                           `json:"$schema"`
	Script      []string                         `json:"script"`
}
