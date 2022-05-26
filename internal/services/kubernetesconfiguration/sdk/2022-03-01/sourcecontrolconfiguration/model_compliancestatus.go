package sourcecontrolconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComplianceStatus struct {
	ComplianceState   *ComplianceStateType `json:"complianceState,omitempty"`
	LastConfigApplied *string              `json:"lastConfigApplied,omitempty"`
	Message           *string              `json:"message,omitempty"`
	MessageLevel      *MessageLevelType    `json:"messageLevel,omitempty"`
}

func (o *ComplianceStatus) GetLastConfigAppliedAsTime() (*time.Time, error) {
	if o.LastConfigApplied == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastConfigApplied, "2006-01-02T15:04:05Z07:00")
}

func (o *ComplianceStatus) SetLastConfigAppliedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastConfigApplied = &formatted
}
