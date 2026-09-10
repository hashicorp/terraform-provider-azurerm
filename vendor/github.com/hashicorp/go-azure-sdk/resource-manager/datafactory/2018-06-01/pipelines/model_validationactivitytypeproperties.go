package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidationActivityTypeProperties struct {
	ChildItems  *interface{}     `json:"childItems,omitempty"`
	Dataset     DatasetReference `json:"dataset"`
	MinimumSize *interface{}     `json:"minimumSize,omitempty"`
	Sleep       *interface{}     `json:"sleep,omitempty"`
	Timeout     *interface{}     `json:"timeout,omitempty"`
}
