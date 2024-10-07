package logicalnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPPool struct {
	End        *string         `json:"end,omitempty"`
	IPPoolType *IPPoolTypeEnum `json:"ipPoolType,omitempty"`
	Info       *IPPoolInfo     `json:"info,omitempty"`
	Name       *string         `json:"name,omitempty"`
	Start      *string         `json:"start,omitempty"`
}
