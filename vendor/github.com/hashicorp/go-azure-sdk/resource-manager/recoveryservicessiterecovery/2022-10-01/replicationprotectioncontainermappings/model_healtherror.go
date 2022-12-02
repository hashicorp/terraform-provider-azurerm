package replicationprotectioncontainermappings

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthError struct {
	CreationTimeUtc              *string                           `json:"creationTimeUtc,omitempty"`
	CustomerResolvability        *HealthErrorCustomerResolvability `json:"customerResolvability,omitempty"`
	EntityId                     *string                           `json:"entityId,omitempty"`
	ErrorCategory                *string                           `json:"errorCategory,omitempty"`
	ErrorCode                    *string                           `json:"errorCode,omitempty"`
	ErrorId                      *string                           `json:"errorId,omitempty"`
	ErrorLevel                   *string                           `json:"errorLevel,omitempty"`
	ErrorMessage                 *string                           `json:"errorMessage,omitempty"`
	ErrorSource                  *string                           `json:"errorSource,omitempty"`
	ErrorType                    *string                           `json:"errorType,omitempty"`
	InnerHealthErrors            *[]InnerHealthError               `json:"innerHealthErrors,omitempty"`
	PossibleCauses               *string                           `json:"possibleCauses,omitempty"`
	RecommendedAction            *string                           `json:"recommendedAction,omitempty"`
	RecoveryProviderErrorMessage *string                           `json:"recoveryProviderErrorMessage,omitempty"`
	SummaryMessage               *string                           `json:"summaryMessage,omitempty"`
}

func (o *HealthError) GetCreationTimeUtcAsTime() (*time.Time, error) {
	if o.CreationTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *HealthError) SetCreationTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTimeUtc = &formatted
}
