package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthErrorDetails struct {
	Code              *string            `json:"code,omitempty"`
	Id                *int64             `json:"id,omitempty"`
	Message           *string            `json:"message,omitempty"`
	MessageParameters *map[string]string `json:"messageParameters,omitempty"`
	PossibleCauses    *string            `json:"possibleCauses,omitempty"`
	RecommendedAction *string            `json:"recommendedAction,omitempty"`
	Severity          *string            `json:"severity,omitempty"`
	Source            *string            `json:"source,omitempty"`
	SummaryMessage    *string            `json:"summaryMessage,omitempty"`
}
