package module

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModuleUpdateParameters struct {
	Location   *string                 `json:"location,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *ModuleUpdateProperties `json:"properties,omitempty"`
	Tags       *map[string]string      `json:"tags,omitempty"`
}
