package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NewRelicMonitorResourceUpdateProperties struct {
	AccountCreationSource     *AccountCreationSource     `json:"accountCreationSource,omitempty"`
	NewRelicAccountProperties *NewRelicAccountProperties `json:"newRelicAccountProperties,omitempty"`
	OrgCreationSource         *OrgCreationSource         `json:"orgCreationSource,omitempty"`
	PlanData                  *PlanData                  `json:"planData,omitempty"`
	UserInfo                  *UserInfo                  `json:"userInfo,omitempty"`
}
