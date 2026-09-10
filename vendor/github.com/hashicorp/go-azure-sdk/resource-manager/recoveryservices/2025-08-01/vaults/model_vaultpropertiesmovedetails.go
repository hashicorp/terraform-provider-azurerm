package vaults

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultPropertiesMoveDetails struct {
	CompletionTimeUtc *string `json:"completionTimeUtc,omitempty"`
	OperationId       *string `json:"operationId,omitempty"`
	SourceResourceId  *string `json:"sourceResourceId,omitempty"`
	StartTimeUtc      *string `json:"startTimeUtc,omitempty"`
	TargetResourceId  *string `json:"targetResourceId,omitempty"`
}

func (o *VaultPropertiesMoveDetails) GetCompletionTimeUtcAsTime() (*time.Time, error) {
	if o.CompletionTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CompletionTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *VaultPropertiesMoveDetails) SetCompletionTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CompletionTimeUtc = &formatted
}

func (o *VaultPropertiesMoveDetails) GetStartTimeUtcAsTime() (*time.Time, error) {
	if o.StartTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *VaultPropertiesMoveDetails) SetStartTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTimeUtc = &formatted
}
