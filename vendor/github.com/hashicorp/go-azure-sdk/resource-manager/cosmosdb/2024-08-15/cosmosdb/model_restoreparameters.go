package cosmosdb

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreParameters struct {
	DatabasesToRestore        *[]DatabaseRestoreResource        `json:"databasesToRestore,omitempty"`
	GremlinDatabasesToRestore *[]GremlinDatabaseRestoreResource `json:"gremlinDatabasesToRestore,omitempty"`
	RestoreMode               *RestoreMode                      `json:"restoreMode,omitempty"`
	RestoreSource             *string                           `json:"restoreSource,omitempty"`
	RestoreTimestampInUtc     *string                           `json:"restoreTimestampInUtc,omitempty"`
	RestoreWithTtlDisabled    *bool                             `json:"restoreWithTtlDisabled,omitempty"`
	TablesToRestore           *[]string                         `json:"tablesToRestore,omitempty"`
}

func (o *RestoreParameters) GetRestoreTimestampInUtcAsTime() (*time.Time, error) {
	if o.RestoreTimestampInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RestoreTimestampInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *RestoreParameters) SetRestoreTimestampInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RestoreTimestampInUtc = &formatted
}
