package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleCreationParameter struct {
	Location   *string                              `json:"location,omitempty"`
	Name       *string                              `json:"name,omitempty"`
	Properties *ScheduleCreationParameterProperties `json:"properties,omitempty"`
	Tags       *map[string]string                   `json:"tags,omitempty"`
}
