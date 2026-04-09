package tags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagCount struct {
	Type  *string `json:"type,omitempty"`
	Value *int64  `json:"value,omitempty"`
}
