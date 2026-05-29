package budgets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CurrentSpend struct {
	Amount *float64 `json:"amount,omitempty"`
	Unit   *string  `json:"unit,omitempty"`
}
