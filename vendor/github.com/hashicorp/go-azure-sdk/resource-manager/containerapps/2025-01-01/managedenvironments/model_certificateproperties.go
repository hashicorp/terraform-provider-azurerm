package managedenvironments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProperties struct {
	CertificateKeyVaultProperties *CertificateKeyVaultProperties `json:"certificateKeyVaultProperties,omitempty"`
	ExpirationDate                *string                        `json:"expirationDate,omitempty"`
	IssueDate                     *string                        `json:"issueDate,omitempty"`
	Issuer                        *string                        `json:"issuer,omitempty"`
	Password                      *string                        `json:"password,omitempty"`
	ProvisioningState             *CertificateProvisioningState  `json:"provisioningState,omitempty"`
	PublicKeyHash                 *string                        `json:"publicKeyHash,omitempty"`
	SubjectAlternativeNames       *[]string                      `json:"subjectAlternativeNames,omitempty"`
	SubjectName                   *string                        `json:"subjectName,omitempty"`
	Thumbprint                    *string                        `json:"thumbprint,omitempty"`
	Valid                         *bool                          `json:"valid,omitempty"`
	Value                         *string                        `json:"value,omitempty"`
}

func (o *CertificateProperties) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *CertificateProperties) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}

func (o *CertificateProperties) GetIssueDateAsTime() (*time.Time, error) {
	if o.IssueDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.IssueDate, "2006-01-02T15:04:05Z07:00")
}

func (o *CertificateProperties) SetIssueDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.IssueDate = &formatted
}
