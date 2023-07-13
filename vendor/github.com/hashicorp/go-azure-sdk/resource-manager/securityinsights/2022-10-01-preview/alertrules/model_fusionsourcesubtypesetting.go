package alertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FusionSourceSubTypeSetting struct {
	Enabled                  bool                        `json:"enabled"`
	SeverityFilters          FusionSubTypeSeverityFilter `json:"severityFilters"`
	SourceSubTypeDisplayName *string                     `json:"sourceSubTypeDisplayName,omitempty"`
	SourceSubTypeName        string                      `json:"sourceSubTypeName"`
}
