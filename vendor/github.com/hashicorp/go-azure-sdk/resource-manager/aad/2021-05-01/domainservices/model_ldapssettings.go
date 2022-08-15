package domainservices

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LdapsSettings struct {
	CertificateNotAfter    *string         `json:"certificateNotAfter,omitempty"`
	CertificateThumbprint  *string         `json:"certificateThumbprint,omitempty"`
	ExternalAccess         *ExternalAccess `json:"externalAccess,omitempty"`
	Ldaps                  *Ldaps          `json:"ldaps,omitempty"`
	PfxCertificate         *string         `json:"pfxCertificate,omitempty"`
	PfxCertificatePassword *string         `json:"pfxCertificatePassword,omitempty"`
	PublicCertificate      *string         `json:"publicCertificate,omitempty"`
}

func (o *LdapsSettings) GetCertificateNotAfterAsTime() (*time.Time, error) {
	if o.CertificateNotAfter == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CertificateNotAfter, "2006-01-02T15:04:05Z07:00")
}

func (o *LdapsSettings) SetCertificateNotAfterAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CertificateNotAfter = &formatted
}
