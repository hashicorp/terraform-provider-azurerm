package alertruletemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FusionTemplateSubTypeSeverityFilter struct {
	IsSupported     bool             `json:"isSupported"`
	SeverityFilters *[]AlertSeverity `json:"severityFilters,omitempty"`
}
