package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TroubleshootingDetails struct {
	Detail             *string                              `json:"detail,omitempty"`
	Id                 *string                              `json:"id,omitempty"`
	ReasonType         *string                              `json:"reasonType,omitempty"`
	RecommendedActions *[]TroubleshootingRecommendedActions `json:"recommendedActions,omitempty"`
	Summary            *string                              `json:"summary,omitempty"`
}
