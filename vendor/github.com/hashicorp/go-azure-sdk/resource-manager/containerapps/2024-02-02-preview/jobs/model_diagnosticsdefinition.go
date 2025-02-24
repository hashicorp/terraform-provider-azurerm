package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticsDefinition struct {
	AnalysisTypes    *[]string                 `json:"analysisTypes,omitempty"`
	Author           *string                   `json:"author,omitempty"`
	Category         *string                   `json:"category,omitempty"`
	Description      *string                   `json:"description,omitempty"`
	Id               *string                   `json:"id,omitempty"`
	Name             *string                   `json:"name,omitempty"`
	Score            *float64                  `json:"score,omitempty"`
	SupportTopicList *[]DiagnosticSupportTopic `json:"supportTopicList,omitempty"`
	Type             *string                   `json:"type,omitempty"`
}
