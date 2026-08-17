package webhook

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookCreateOrUpdateProperties struct {
	ExpiryTime *string                     `json:"expiryTime,omitempty"`
	IsEnabled  *bool                       `json:"isEnabled,omitempty"`
	Parameters *map[string]string          `json:"parameters,omitempty"`
	RunOn      *string                     `json:"runOn,omitempty"`
	Runbook    *RunbookAssociationProperty `json:"runbook,omitempty"`
	Uri        *string                     `json:"uri,omitempty"`
}

func (o *WebhookCreateOrUpdateProperties) GetExpiryTimeAsTime() (*time.Time, error) {
	if o.ExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WebhookCreateOrUpdateProperties) SetExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTime = &formatted
}
