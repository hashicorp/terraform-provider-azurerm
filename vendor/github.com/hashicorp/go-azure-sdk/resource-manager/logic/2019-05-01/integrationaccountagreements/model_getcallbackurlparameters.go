package integrationaccountagreements

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetCallbackURLParameters struct {
	KeyType  *KeyType `json:"keyType,omitempty"`
	NotAfter *string  `json:"notAfter,omitempty"`
}

func (o *GetCallbackURLParameters) GetNotAfterAsTime() (*time.Time, error) {
	if o.NotAfter == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NotAfter, "2006-01-02T15:04:05Z07:00")
}

func (o *GetCallbackURLParameters) SetNotAfterAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NotAfter = &formatted
}
