package appserviceenvironments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SlotSwapStatus struct {
	DestinationSlotName *string `json:"destinationSlotName,omitempty"`
	SourceSlotName      *string `json:"sourceSlotName,omitempty"`
	TimestampUtc        *string `json:"timestampUtc,omitempty"`
}

func (o *SlotSwapStatus) GetTimestampUtcAsTime() (*time.Time, error) {
	if o.TimestampUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimestampUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *SlotSwapStatus) SetTimestampUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimestampUtc = &formatted
}
