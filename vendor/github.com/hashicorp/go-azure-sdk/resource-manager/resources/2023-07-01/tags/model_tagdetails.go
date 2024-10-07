package tags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagDetails struct {
	Count   *TagCount   `json:"count,omitempty"`
	Id      *string     `json:"id,omitempty"`
	TagName *string     `json:"tagName,omitempty"`
	Values  *[]TagValue `json:"values,omitempty"`
}
