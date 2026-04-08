package domains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JsonFieldWithDefault struct {
	DefaultValue *string `json:"defaultValue,omitempty"`
	SourceField  *string `json:"sourceField,omitempty"`
}
