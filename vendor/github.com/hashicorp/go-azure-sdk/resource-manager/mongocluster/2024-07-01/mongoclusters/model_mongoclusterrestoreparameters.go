package mongoclusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoClusterRestoreParameters struct {
	PointInTimeUTC   *string `json:"pointInTimeUTC,omitempty"`
	SourceResourceId *string `json:"sourceResourceId,omitempty"`
}

func (o *MongoClusterRestoreParameters) GetPointInTimeUTCAsTime() (*time.Time, error) {
	if o.PointInTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PointInTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *MongoClusterRestoreParameters) SetPointInTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PointInTimeUTC = &formatted
}
