package automations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationTriggeringRule struct {
	ExpectedValue *string       `json:"expectedValue,omitempty"`
	Operator      *Operator     `json:"operator,omitempty"`
	PropertyJPath *string       `json:"propertyJPath,omitempty"`
	PropertyType  *PropertyType `json:"propertyType,omitempty"`
}
