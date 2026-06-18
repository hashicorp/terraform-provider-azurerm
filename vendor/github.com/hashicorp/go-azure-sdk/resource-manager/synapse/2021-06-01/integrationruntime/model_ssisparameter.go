package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SsisParameter struct {
	DataType              *string `json:"dataType,omitempty"`
	DefaultValue          *string `json:"defaultValue,omitempty"`
	Description           *string `json:"description,omitempty"`
	DesignDefaultValue    *string `json:"designDefaultValue,omitempty"`
	Id                    *int64  `json:"id,omitempty"`
	Name                  *string `json:"name,omitempty"`
	Required              *bool   `json:"required,omitempty"`
	Sensitive             *bool   `json:"sensitive,omitempty"`
	SensitiveDefaultValue *string `json:"sensitiveDefaultValue,omitempty"`
	ValueSet              *bool   `json:"valueSet,omitempty"`
	ValueType             *string `json:"valueType,omitempty"`
	Variable              *string `json:"variable,omitempty"`
}
