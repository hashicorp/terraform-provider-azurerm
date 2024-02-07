package policy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyCollection struct {
	Count    *int64            `json:"count,omitempty"`
	NextLink *string           `json:"nextLink,omitempty"`
	Value    *[]PolicyContract `json:"value,omitempty"`
}
