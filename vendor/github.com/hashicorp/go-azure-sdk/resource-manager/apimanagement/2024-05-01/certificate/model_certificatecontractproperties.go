package certificate

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateContractProperties struct {
	ExpirationDate string                      `json:"expirationDate"`
	KeyVault       *KeyVaultContractProperties `json:"keyVault,omitempty"`
	Subject        string                      `json:"subject"`
	Thumbprint     string                      `json:"thumbprint"`
}

func (o *CertificateContractProperties) GetExpirationDateAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *CertificateContractProperties) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = formatted
}
