package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApplication struct {
	ApplicationPool *string `json:"applicationPool,omitempty"`
	GroupName       *string `json:"groupName,omitempty"`
	Name            *string `json:"name,omitempty"`
	Platform        *string `json:"platform,omitempty"`
	Status          *string `json:"status,omitempty"`
	WebServer       *string `json:"webServer,omitempty"`
}
