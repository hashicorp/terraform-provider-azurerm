package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidationActivityTypeProperties struct {
	ChildItems  *bool            `json:"childItems,omitempty"`
	Dataset     DatasetReference `json:"dataset"`
	MinimumSize *int64           `json:"minimumSize,omitempty"`
	Sleep       *int64           `json:"sleep,omitempty"`
	Timeout     *string          `json:"timeout,omitempty"`
}
