package cosmosdb

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreParametersBase struct {
	RestoreSource          *string `json:"restoreSource,omitempty"`
	RestoreTimestampInUtc  *string `json:"restoreTimestampInUtc,omitempty"`
	RestoreWithTtlDisabled *bool   `json:"restoreWithTtlDisabled,omitempty"`
}

func (o *RestoreParametersBase) GetRestoreTimestampInUtcAsTime() (*time.Time, error) {
	if o.RestoreTimestampInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RestoreTimestampInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *RestoreParametersBase) SetRestoreTimestampInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RestoreTimestampInUtc = &formatted
}
