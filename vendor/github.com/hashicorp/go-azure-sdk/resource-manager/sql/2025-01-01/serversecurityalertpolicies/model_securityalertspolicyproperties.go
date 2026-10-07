package serversecurityalertpolicies

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityAlertsPolicyProperties struct {
	CreationTime            *string                   `json:"creationTime,omitempty"`
	DisabledAlerts          *[]string                 `json:"disabledAlerts,omitempty"`
	EmailAccountAdmins      *bool                     `json:"emailAccountAdmins,omitempty"`
	EmailAddresses          *[]string                 `json:"emailAddresses,omitempty"`
	RetentionDays           *int64                    `json:"retentionDays,omitempty"`
	State                   SecurityAlertsPolicyState `json:"state"`
	StorageAccountAccessKey *string                   `json:"storageAccountAccessKey,omitempty"`
	StorageEndpoint         *string                   `json:"storageEndpoint,omitempty"`
}

func (o *SecurityAlertsPolicyProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SecurityAlertsPolicyProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}
