package actionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchObject struct {
	Properties *PatchProperties `json:"properties,omitempty"`
	Tags       *interface{}     `json:"tags,omitempty"`
}
