package factories

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAccessPolicy struct {
	AccessResourcePath *string `json:"accessResourcePath,omitempty"`
	ExpireTime         *string `json:"expireTime,omitempty"`
	Permissions        *string `json:"permissions,omitempty"`
	ProfileName        *string `json:"profileName,omitempty"`
	StartTime          *string `json:"startTime,omitempty"`
}
