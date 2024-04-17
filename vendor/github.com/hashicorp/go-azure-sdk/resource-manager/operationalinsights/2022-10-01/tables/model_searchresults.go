package tables

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchResults struct {
	AzureAsyncOperationId *string `json:"azureAsyncOperationId,omitempty"`
	Description           *string `json:"description,omitempty"`
	EndSearchTime         *string `json:"endSearchTime,omitempty"`
	Limit                 *int64  `json:"limit,omitempty"`
	Query                 *string `json:"query,omitempty"`
	SourceTable           *string `json:"sourceTable,omitempty"`
	StartSearchTime       *string `json:"startSearchTime,omitempty"`
}

func (o *SearchResults) GetEndSearchTimeAsTime() (*time.Time, error) {
	if o.EndSearchTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndSearchTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SearchResults) SetEndSearchTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndSearchTime = &formatted
}

func (o *SearchResults) GetStartSearchTimeAsTime() (*time.Time, error) {
	if o.StartSearchTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartSearchTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SearchResults) SetStartSearchTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartSearchTime = &formatted
}
