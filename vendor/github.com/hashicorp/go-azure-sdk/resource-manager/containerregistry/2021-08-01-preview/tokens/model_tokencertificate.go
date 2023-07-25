package tokens

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenCertificate struct {
	EncodedPemCertificate *string               `json:"encodedPemCertificate,omitempty"`
	Expiry                *string               `json:"expiry,omitempty"`
	Name                  *TokenCertificateName `json:"name,omitempty"`
	Thumbprint            *string               `json:"thumbprint,omitempty"`
}

func (o *TokenCertificate) GetExpiryAsTime() (*time.Time, error) {
	if o.Expiry == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Expiry, "2006-01-02T15:04:05Z07:00")
}

func (o *TokenCertificate) SetExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Expiry = &formatted
}
