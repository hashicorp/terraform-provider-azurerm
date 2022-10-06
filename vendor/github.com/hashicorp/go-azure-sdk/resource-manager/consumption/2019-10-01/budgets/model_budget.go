package budgets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Budget struct {
	ETag       *string           `json:"eTag,omitempty"`
	Id         *string           `json:"id,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties *BudgetProperties `json:"properties,omitempty"`
	Type       *string           `json:"type,omitempty"`
}
