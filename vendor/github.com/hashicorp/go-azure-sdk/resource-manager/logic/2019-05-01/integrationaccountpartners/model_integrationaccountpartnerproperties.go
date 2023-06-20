package integrationaccountpartners

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountPartnerProperties struct {
	ChangedTime *string        `json:"changedTime,omitempty"`
	Content     PartnerContent `json:"content"`
	CreatedTime *string        `json:"createdTime,omitempty"`
	Metadata    *interface{}   `json:"metadata,omitempty"`
	PartnerType PartnerType    `json:"partnerType"`
}

func (o *IntegrationAccountPartnerProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountPartnerProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *IntegrationAccountPartnerProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountPartnerProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
