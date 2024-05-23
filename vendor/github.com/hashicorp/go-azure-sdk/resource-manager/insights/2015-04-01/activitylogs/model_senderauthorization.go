package activitylogs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SenderAuthorization struct {
	Action *string `json:"action,omitempty"`
	Role   *string `json:"role,omitempty"`
	Scope  *string `json:"scope,omitempty"`
}
