package replicationfabrics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionDetails struct {
	KekCertExpiryDate *string `json:"kekCertExpiryDate,omitempty"`
	KekCertThumbprint *string `json:"kekCertThumbprint,omitempty"`
	KekState          *string `json:"kekState,omitempty"`
}

func (o *EncryptionDetails) GetKekCertExpiryDateAsTime() (*time.Time, error) {
	if o.KekCertExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.KekCertExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *EncryptionDetails) SetKekCertExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.KekCertExpiryDate = &formatted
}
