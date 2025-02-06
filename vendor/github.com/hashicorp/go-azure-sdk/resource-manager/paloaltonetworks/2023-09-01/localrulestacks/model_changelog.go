package localrulestacks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Changelog struct {
	Changes       []string `json:"changes"`
	LastCommitted *string  `json:"lastCommitted,omitempty"`
	LastModified  *string  `json:"lastModified,omitempty"`
}

func (o *Changelog) GetLastCommittedAsTime() (*time.Time, error) {
	if o.LastCommitted == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCommitted, "2006-01-02T15:04:05Z07:00")
}

func (o *Changelog) SetLastCommittedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCommitted = &formatted
}

func (o *Changelog) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *Changelog) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
