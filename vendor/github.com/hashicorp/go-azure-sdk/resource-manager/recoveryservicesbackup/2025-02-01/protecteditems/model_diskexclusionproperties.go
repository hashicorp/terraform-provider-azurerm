package protecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskExclusionProperties struct {
	DiskLunList     *[]int64 `json:"diskLunList,omitempty"`
	IsInclusionList *bool    `json:"isInclusionList,omitempty"`
}
