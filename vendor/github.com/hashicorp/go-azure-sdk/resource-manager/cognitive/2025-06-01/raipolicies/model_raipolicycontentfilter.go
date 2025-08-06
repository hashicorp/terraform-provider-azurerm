package raipolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RaiPolicyContentFilter struct {
	Blocking          *bool                   `json:"blocking,omitempty"`
	Enabled           *bool                   `json:"enabled,omitempty"`
	Name              *string                 `json:"name,omitempty"`
	SeverityThreshold *ContentLevel           `json:"severityThreshold,omitempty"`
	Source            *RaiPolicyContentSource `json:"source,omitempty"`
}
