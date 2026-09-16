package buckets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketServerProperties struct {
	CertificateCommonName       *string                      `json:"certificateCommonName,omitempty"`
	CertificateExpiryDate       *string                      `json:"certificateExpiryDate,omitempty"`
	CertificateObject           *string                      `json:"certificateObject,omitempty"`
	Fqdn                        *string                      `json:"fqdn,omitempty"`
	IPAddress                   *string                      `json:"ipAddress,omitempty"`
	OnCertificateConflictAction *OnCertificateConflictAction `json:"onCertificateConflictAction,omitempty"`
}

func (o *BucketServerProperties) GetCertificateExpiryDateAsTime() (*time.Time, error) {
	if o.CertificateExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CertificateExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BucketServerProperties) SetCertificateExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CertificateExpiryDate = &formatted
}
