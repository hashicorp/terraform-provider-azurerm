package firewallpolicykubeselectorgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabelSelectorExpression struct {
	Key      *string                `json:"key,omitempty"`
	Operator *LabelSelectorOperator `json:"operator,omitempty"`
	Values   *[]string              `json:"values,omitempty"`
}
