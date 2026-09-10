package webhook

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookProperties struct {
	CreationTime     *string                     `json:"creationTime,omitempty"`
	Description      *string                     `json:"description,omitempty"`
	ExpiryTime       *string                     `json:"expiryTime,omitempty"`
	IsEnabled        *bool                       `json:"isEnabled,omitempty"`
	LastInvokedTime  *string                     `json:"lastInvokedTime,omitempty"`
	LastModifiedBy   *string                     `json:"lastModifiedBy,omitempty"`
	LastModifiedTime *string                     `json:"lastModifiedTime,omitempty"`
	Parameters       *map[string]string          `json:"parameters,omitempty"`
	RunOn            *string                     `json:"runOn,omitempty"`
	Runbook          *RunbookAssociationProperty `json:"runbook,omitempty"`
	Uri              *string                     `json:"uri,omitempty"`
}

func (o *WebhookProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WebhookProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *WebhookProperties) GetExpiryTimeAsTime() (*time.Time, error) {
	if o.ExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WebhookProperties) SetExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTime = &formatted
}

func (o *WebhookProperties) GetLastInvokedTimeAsTime() (*time.Time, error) {
	if o.LastInvokedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastInvokedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WebhookProperties) SetLastInvokedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastInvokedTime = &formatted
}

func (o *WebhookProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WebhookProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
