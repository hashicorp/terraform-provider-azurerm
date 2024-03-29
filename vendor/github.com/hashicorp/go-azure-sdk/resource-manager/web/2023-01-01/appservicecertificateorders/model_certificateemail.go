package appservicecertificateorders

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateEmail struct {
	EmailId   *string `json:"emailId,omitempty"`
	TimeStamp *string `json:"timeStamp,omitempty"`
}

func (o *CertificateEmail) GetTimeStampAsTime() (*time.Time, error) {
	if o.TimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CertificateEmail) SetTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeStamp = &formatted
}
