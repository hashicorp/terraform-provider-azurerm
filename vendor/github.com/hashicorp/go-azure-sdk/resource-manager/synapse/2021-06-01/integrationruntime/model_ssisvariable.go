package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SsisVariable struct {
	DataType       *string `json:"dataType,omitempty"`
	Description    *string `json:"description,omitempty"`
	Id             *int64  `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Sensitive      *bool   `json:"sensitive,omitempty"`
	SensitiveValue *string `json:"sensitiveValue,omitempty"`
	Value          *string `json:"value,omitempty"`
}
