package namedvalue

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultLastAccessStatusContractProperties struct {
	Code         *string `json:"code,omitempty"`
	Message      *string `json:"message,omitempty"`
	TimeStampUtc *string `json:"timeStampUtc,omitempty"`
}

func (o *KeyVaultLastAccessStatusContractProperties) GetTimeStampUtcAsTime() (*time.Time, error) {
	if o.TimeStampUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeStampUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyVaultLastAccessStatusContractProperties) SetTimeStampUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeStampUtc = &formatted
}
