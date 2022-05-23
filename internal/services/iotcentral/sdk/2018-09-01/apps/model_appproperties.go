package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppProperties struct {
	ApplicationId *string `json:"applicationId,omitempty"`
	DisplayName   *string `json:"displayName,omitempty"`
	Subdomain     *string `json:"subdomain,omitempty"`
	Template      *string `json:"template,omitempty"`
}
