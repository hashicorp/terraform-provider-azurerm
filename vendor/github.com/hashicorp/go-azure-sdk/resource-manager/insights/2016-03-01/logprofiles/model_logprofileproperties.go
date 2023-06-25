package logprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogProfileProperties struct {
	Categories       []string        `json:"categories"`
	Locations        []string        `json:"locations"`
	RetentionPolicy  RetentionPolicy `json:"retentionPolicy"`
	ServiceBusRuleId *string         `json:"serviceBusRuleId,omitempty"`
	StorageAccountId *string         `json:"storageAccountId,omitempty"`
}
