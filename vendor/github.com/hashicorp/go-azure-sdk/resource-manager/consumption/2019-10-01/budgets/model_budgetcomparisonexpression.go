package budgets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BudgetComparisonExpression struct {
	Name     string             `json:"name"`
	Operator BudgetOperatorType `json:"operator"`
	Values   []string           `json:"values"`
}
