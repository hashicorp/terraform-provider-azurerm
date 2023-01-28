package integrationaccountagreements

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountAgreementProperties struct {
	AgreementType AgreementType    `json:"agreementType"`
	ChangedTime   *string          `json:"changedTime,omitempty"`
	Content       AgreementContent `json:"content"`
	CreatedTime   *string          `json:"createdTime,omitempty"`
	GuestIdentity BusinessIdentity `json:"guestIdentity"`
	GuestPartner  string           `json:"guestPartner"`
	HostIdentity  BusinessIdentity `json:"hostIdentity"`
	HostPartner   string           `json:"hostPartner"`
	Metadata      *interface{}     `json:"metadata,omitempty"`
}

func (o *IntegrationAccountAgreementProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountAgreementProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *IntegrationAccountAgreementProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountAgreementProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
