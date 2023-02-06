package updatesummaries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PackageVersionInfo struct {
	LastUpdated *string `json:"lastUpdated,omitempty"`
	PackageType *string `json:"packageType,omitempty"`
	Version     *string `json:"version,omitempty"`
}

func (o *PackageVersionInfo) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *PackageVersionInfo) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}
