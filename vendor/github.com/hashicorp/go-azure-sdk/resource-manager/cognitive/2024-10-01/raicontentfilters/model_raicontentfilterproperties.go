package raicontentfilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RaiContentFilterProperties struct {
	IsMultiLevelFilter *bool                   `json:"isMultiLevelFilter,omitempty"`
	Name               *string                 `json:"name,omitempty"`
	Source             *RaiPolicyContentSource `json:"source,omitempty"`
}
