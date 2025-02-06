package projects

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectUpdate struct {
	Location   *string                  `json:"location,omitempty"`
	Properties *ProjectUpdateProperties `json:"properties,omitempty"`
	Tags       *map[string]string       `json:"tags,omitempty"`
}
