package apimanagementservice

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateInformation struct {
	Expiry     string `json:"expiry"`
	Subject    string `json:"subject"`
	Thumbprint string `json:"thumbprint"`
}

func (o *CertificateInformation) GetExpiryAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.Expiry, "2006-01-02T15:04:05Z07:00")
}

func (o *CertificateInformation) SetExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Expiry = formatted
}
