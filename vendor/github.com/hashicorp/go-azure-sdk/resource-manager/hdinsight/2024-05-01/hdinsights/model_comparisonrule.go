package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComparisonRule struct {
	Operator  ComparisonOperator `json:"operator"`
	Threshold float64            `json:"threshold"`
}
