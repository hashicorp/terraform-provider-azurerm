package managedenvironments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainConfiguration struct {
	CertificateKeyVaultProperties *CertificateKeyVaultProperties `json:"certificateKeyVaultProperties,omitempty"`
	CertificatePassword           *string                        `json:"certificatePassword,omitempty"`
	CertificateValue              *string                        `json:"certificateValue,omitempty"`
	CustomDomainVerificationId    *string                        `json:"customDomainVerificationId,omitempty"`
	DnsSuffix                     *string                        `json:"dnsSuffix,omitempty"`
	ExpirationDate                *string                        `json:"expirationDate,omitempty"`
	SubjectName                   *string                        `json:"subjectName,omitempty"`
	Thumbprint                    *string                        `json:"thumbprint,omitempty"`
}

func (o *CustomDomainConfiguration) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *CustomDomainConfiguration) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}
