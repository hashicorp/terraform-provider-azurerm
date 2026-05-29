package alertprocessingrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Condition struct {
	Field    *Field    `json:"field,omitempty"`
	Operator *Operator `json:"operator,omitempty"`
	Values   *[]string `json:"values,omitempty"`
}
