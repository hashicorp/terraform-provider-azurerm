package cacertificates

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CaCertificateProperties struct {
	Description        *string                         `json:"description,omitempty"`
	EncodedCertificate *string                         `json:"encodedCertificate,omitempty"`
	ExpiryTimeInUtc    *string                         `json:"expiryTimeInUtc,omitempty"`
	IssueTimeInUtc     *string                         `json:"issueTimeInUtc,omitempty"`
	ProvisioningState  *CaCertificateProvisioningState `json:"provisioningState,omitempty"`
}

func (o *CaCertificateProperties) GetExpiryTimeInUtcAsTime() (*time.Time, error) {
	if o.ExpiryTimeInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTimeInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *CaCertificateProperties) SetExpiryTimeInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTimeInUtc = &formatted
}

func (o *CaCertificateProperties) GetIssueTimeInUtcAsTime() (*time.Time, error) {
	if o.IssueTimeInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.IssueTimeInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *CaCertificateProperties) SetIssueTimeInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.IssueTimeInUtc = &formatted
}
