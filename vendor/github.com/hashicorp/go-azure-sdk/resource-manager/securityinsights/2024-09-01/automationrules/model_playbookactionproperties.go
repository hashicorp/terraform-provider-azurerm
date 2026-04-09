package automationrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlaybookActionProperties struct {
	LogicAppResourceId string  `json:"logicAppResourceId"`
	TenantId           *string `json:"tenantId,omitempty"`
}
