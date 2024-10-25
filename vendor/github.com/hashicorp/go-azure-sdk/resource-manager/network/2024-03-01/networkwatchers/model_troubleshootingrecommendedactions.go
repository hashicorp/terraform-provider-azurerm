package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TroubleshootingRecommendedActions struct {
	ActionId      *string `json:"actionId,omitempty"`
	ActionText    *string `json:"actionText,omitempty"`
	ActionUri     *string `json:"actionUri,omitempty"`
	ActionUriText *string `json:"actionUriText,omitempty"`
}
