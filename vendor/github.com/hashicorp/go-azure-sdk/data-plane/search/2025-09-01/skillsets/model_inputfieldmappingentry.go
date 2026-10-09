package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputFieldMappingEntry struct {
	Inputs        *[]InputFieldMappingEntry `json:"inputs,omitempty"`
	Name          string                    `json:"name"`
	Source        *string                   `json:"source,omitempty"`
	SourceContext *string                   `json:"sourceContext,omitempty"`
}
