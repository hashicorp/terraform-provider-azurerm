package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceSasParameters struct {
	CanonicalizedResource string          `json:"canonicalizedResource"`
	EndPk                 *string         `json:"endPk,omitempty"`
	EndRk                 *string         `json:"endRk,omitempty"`
	KeyToSign             *string         `json:"keyToSign,omitempty"`
	Rscc                  *string         `json:"rscc,omitempty"`
	Rscd                  *string         `json:"rscd,omitempty"`
	Rsce                  *string         `json:"rsce,omitempty"`
	Rscl                  *string         `json:"rscl,omitempty"`
	Rsct                  *string         `json:"rsct,omitempty"`
	SignedExpiry          *string         `json:"signedExpiry,omitempty"`
	SignedIP              *string         `json:"signedIp,omitempty"`
	SignedIdentifier      *string         `json:"signedIdentifier,omitempty"`
	SignedPermission      *Permissions    `json:"signedPermission,omitempty"`
	SignedProtocol        *HTTPProtocol   `json:"signedProtocol,omitempty"`
	SignedResource        *SignedResource `json:"signedResource,omitempty"`
	SignedStart           *string         `json:"signedStart,omitempty"`
	StartPk               *string         `json:"startPk,omitempty"`
	StartRk               *string         `json:"startRk,omitempty"`
}

func (o *ServiceSasParameters) GetSignedExpiryAsTime() (*time.Time, error) {
	if o.SignedExpiry == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SignedExpiry, "2006-01-02T15:04:05Z07:00")
}

func (o *ServiceSasParameters) SetSignedExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SignedExpiry = &formatted
}

func (o *ServiceSasParameters) GetSignedStartAsTime() (*time.Time, error) {
	if o.SignedStart == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SignedStart, "2006-01-02T15:04:05Z07:00")
}

func (o *ServiceSasParameters) SetSignedStartAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SignedStart = &formatted
}
