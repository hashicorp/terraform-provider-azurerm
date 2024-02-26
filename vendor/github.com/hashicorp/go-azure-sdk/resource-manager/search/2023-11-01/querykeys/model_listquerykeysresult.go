package querykeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListQueryKeysResult struct {
	NextLink *string     `json:"nextLink,omitempty"`
	Value    *[]QueryKey `json:"value,omitempty"`
}
