package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestNotificationDetailsResponse struct {
	ActionDetails *[]ActionDetail `json:"actionDetails,omitempty"`
	CompletedTime *string         `json:"completedTime,omitempty"`
	Context       *Context        `json:"context,omitempty"`
	CreatedTime   *string         `json:"createdTime,omitempty"`
	State         string          `json:"state"`
}
