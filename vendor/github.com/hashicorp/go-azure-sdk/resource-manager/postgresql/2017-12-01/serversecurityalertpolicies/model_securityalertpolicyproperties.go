package serversecurityalertpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityAlertPolicyProperties struct {
	DisabledAlerts          *[]string                      `json:"disabledAlerts,omitempty"`
	EmailAccountAdmins      *bool                          `json:"emailAccountAdmins,omitempty"`
	EmailAddresses          *[]string                      `json:"emailAddresses,omitempty"`
	RetentionDays           *int64                         `json:"retentionDays,omitempty"`
	State                   ServerSecurityAlertPolicyState `json:"state"`
	StorageAccountAccessKey *string                        `json:"storageAccountAccessKey,omitempty"`
	StorageEndpoint         *string                        `json:"storageEndpoint,omitempty"`
}
