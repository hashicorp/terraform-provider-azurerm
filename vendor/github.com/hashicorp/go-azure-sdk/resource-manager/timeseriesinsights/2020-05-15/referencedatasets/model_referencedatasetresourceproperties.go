package referencedatasets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReferenceDataSetResourceProperties struct {
	CreationTime                 *string                       `json:"creationTime,omitempty"`
	DataStringComparisonBehavior *DataStringComparisonBehavior `json:"dataStringComparisonBehavior,omitempty"`
	KeyProperties                []ReferenceDataSetKeyProperty `json:"keyProperties"`
	ProvisioningState            *ProvisioningState            `json:"provisioningState,omitempty"`
}

func (o *ReferenceDataSetResourceProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ReferenceDataSetResourceProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}
