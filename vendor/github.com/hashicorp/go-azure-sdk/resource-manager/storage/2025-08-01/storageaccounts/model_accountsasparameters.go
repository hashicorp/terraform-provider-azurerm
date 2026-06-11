package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountSasParameters struct {
	KeyToSign           *string             `json:"keyToSign,omitempty"`
	SignedExpiry        string              `json:"signedExpiry"`
	SignedIP            *string             `json:"signedIp,omitempty"`
	SignedPermission    Permissions         `json:"signedPermission"`
	SignedProtocol      *HTTPProtocol       `json:"signedProtocol,omitempty"`
	SignedResourceTypes SignedResourceTypes `json:"signedResourceTypes"`
	SignedServices      Services            `json:"signedServices"`
	SignedStart         *string             `json:"signedStart,omitempty"`
}

func (o *AccountSasParameters) GetSignedExpiryAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.SignedExpiry, "2006-01-02T15:04:05Z07:00")
}

func (o *AccountSasParameters) SetSignedExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SignedExpiry = formatted
}

func (o *AccountSasParameters) GetSignedStartAsTime() (*time.Time, error) {
	if o.SignedStart == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SignedStart, "2006-01-02T15:04:05Z07:00")
}

func (o *AccountSasParameters) SetSignedStartAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SignedStart = &formatted
}
