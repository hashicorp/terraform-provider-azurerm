package nodepools

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Condition struct {
	LastTransitionTime string        `json:"lastTransitionTime"`
	Message            string        `json:"message"`
	Reason             string        `json:"reason"`
	Status             StatusType    `json:"status"`
	Type               ConditionType `json:"type"`
}

func (o *Condition) GetLastTransitionTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.LastTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Condition) SetLastTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTransitionTime = formatted
}
