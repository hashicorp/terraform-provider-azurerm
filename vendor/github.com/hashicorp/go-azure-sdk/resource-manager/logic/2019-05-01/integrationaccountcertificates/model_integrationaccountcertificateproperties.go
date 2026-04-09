package integrationaccountcertificates

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountCertificateProperties struct {
	ChangedTime       *string               `json:"changedTime,omitempty"`
	CreatedTime       *string               `json:"createdTime,omitempty"`
	Key               *KeyVaultKeyReference `json:"key,omitempty"`
	Metadata          *interface{}          `json:"metadata,omitempty"`
	PublicCertificate *string               `json:"publicCertificate,omitempty"`
}

func (o *IntegrationAccountCertificateProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountCertificateProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *IntegrationAccountCertificateProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountCertificateProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
