package protectioncontainers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DPMContainerExtendedInfo struct {
	LastRefreshedAt *string `json:"lastRefreshedAt,omitempty"`
}

func (o *DPMContainerExtendedInfo) GetLastRefreshedAtAsTime() (*time.Time, error) {
	if o.LastRefreshedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRefreshedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *DPMContainerExtendedInfo) SetLastRefreshedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRefreshedAt = &formatted
}
