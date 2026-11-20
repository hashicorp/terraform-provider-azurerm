package edgeactionexecutionfilters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionExecutionFilterProperties struct {
	ExecutionFilterIdentifierHeaderName  string             `json:"executionFilterIdentifierHeaderName"`
	ExecutionFilterIdentifierHeaderValue string             `json:"executionFilterIdentifierHeaderValue"`
	LastUpdateTime                       *string            `json:"lastUpdateTime,omitempty"`
	ProvisioningState                    *ProvisioningState `json:"provisioningState,omitempty"`
	VersionId                            string             `json:"versionId"`
}

func (o *EdgeActionExecutionFilterProperties) GetLastUpdateTimeAsTime() (*time.Time, error) {
	if o.LastUpdateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EdgeActionExecutionFilterProperties) SetLastUpdateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdateTime = &formatted
}
