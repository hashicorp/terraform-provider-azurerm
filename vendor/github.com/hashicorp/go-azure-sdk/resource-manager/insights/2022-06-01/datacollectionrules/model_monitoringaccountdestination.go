package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoringAccountDestination struct {
	AccountId         *string `json:"accountId,omitempty"`
	AccountResourceId *string `json:"accountResourceId,omitempty"`
	Name              *string `json:"name,omitempty"`
}
