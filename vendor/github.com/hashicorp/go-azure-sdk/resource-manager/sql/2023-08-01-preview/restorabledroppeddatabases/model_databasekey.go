package restorabledroppeddatabases

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseKey struct {
	CreationDate *string          `json:"creationDate,omitempty"`
	Subregion    *string          `json:"subregion,omitempty"`
	Thumbprint   *string          `json:"thumbprint,omitempty"`
	Type         *DatabaseKeyType `json:"type,omitempty"`
}

func (o *DatabaseKey) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *DatabaseKey) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
