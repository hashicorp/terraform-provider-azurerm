package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SQLInstanceSettings struct {
	Collation                          *string `json:"collation,omitempty"`
	IsIfiEnabled                       *bool   `json:"isIfiEnabled,omitempty"`
	IsLpimEnabled                      *bool   `json:"isLpimEnabled,omitempty"`
	IsOptimizeForAdHocWorkloadsEnabled *bool   `json:"isOptimizeForAdHocWorkloadsEnabled,omitempty"`
	MaxDop                             *int64  `json:"maxDop,omitempty"`
	MaxServerMemoryMB                  *int64  `json:"maxServerMemoryMB,omitempty"`
	MinServerMemoryMB                  *int64  `json:"minServerMemoryMB,omitempty"`
}
