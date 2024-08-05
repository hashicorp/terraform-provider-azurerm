package tags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagValue struct {
	Count    *TagCount `json:"count,omitempty"`
	Id       *string   `json:"id,omitempty"`
	TagValue *string   `json:"tagValue,omitempty"`
}
