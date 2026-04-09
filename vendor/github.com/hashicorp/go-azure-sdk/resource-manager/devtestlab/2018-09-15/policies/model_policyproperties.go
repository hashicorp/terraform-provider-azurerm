package policies

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyProperties struct {
	CreatedDate       *string              `json:"createdDate,omitempty"`
	Description       *string              `json:"description,omitempty"`
	EvaluatorType     *PolicyEvaluatorType `json:"evaluatorType,omitempty"`
	FactData          *string              `json:"factData,omitempty"`
	FactName          *PolicyFactName      `json:"factName,omitempty"`
	ProvisioningState *string              `json:"provisioningState,omitempty"`
	Status            *PolicyStatus        `json:"status,omitempty"`
	Threshold         *string              `json:"threshold,omitempty"`
	UniqueIdentifier  *string              `json:"uniqueIdentifier,omitempty"`
}

func (o *PolicyProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *PolicyProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}
