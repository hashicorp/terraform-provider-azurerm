package autonomousdatabases

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreAutonomousDatabaseDetails struct {
	Timestamp string `json:"timestamp"`
}

func (o *RestoreAutonomousDatabaseDetails) GetTimestampAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *RestoreAutonomousDatabaseDetails) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = formatted
}
