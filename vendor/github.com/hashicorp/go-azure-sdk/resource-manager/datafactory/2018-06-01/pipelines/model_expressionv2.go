package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressionV2 struct {
	Operands  *[]ExpressionV2   `json:"operands,omitempty"`
	Operators *[]string         `json:"operators,omitempty"`
	Type      *ExpressionV2Type `json:"type,omitempty"`
	Value     *interface{}      `json:"value,omitempty"`
}
