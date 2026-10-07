package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportImageParameters struct {
	Mode                       *ImportMode  `json:"mode,omitempty"`
	Source                     ImportSource `json:"source"`
	TargetTags                 *[]string    `json:"targetTags,omitempty"`
	UntaggedTargetRepositories *[]string    `json:"untaggedTargetRepositories,omitempty"`
}
