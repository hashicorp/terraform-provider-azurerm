package serviceresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateMySqlAzureDbForMySqlOfflineTaskInput struct {
	MakeSourceServerReadOnly *bool                                             `json:"makeSourceServerReadOnly,omitempty"`
	OptionalAgentSettings    *map[string]string                                `json:"optionalAgentSettings,omitempty"`
	SelectedDatabases        []MigrateMySqlAzureDbForMySqlOfflineDatabaseInput `json:"selectedDatabases"`
	SourceConnectionInfo     MySqlConnectionInfo                               `json:"sourceConnectionInfo"`
	StartedOn                *string                                           `json:"startedOn,omitempty"`
	TargetConnectionInfo     MySqlConnectionInfo                               `json:"targetConnectionInfo"`
}

func (o *MigrateMySqlAzureDbForMySqlOfflineTaskInput) GetStartedOnAsTime() (*time.Time, error) {
	if o.StartedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *MigrateMySqlAzureDbForMySqlOfflineTaskInput) SetStartedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedOn = &formatted
}
