package assessments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentStatusResponse struct {
	Cause               *string              `json:"cause,omitempty"`
	Code                AssessmentStatusCode `json:"code"`
	Description         *string              `json:"description,omitempty"`
	FirstEvaluationDate *string              `json:"firstEvaluationDate,omitempty"`
	StatusChangeDate    *string              `json:"statusChangeDate,omitempty"`
}

func (o *AssessmentStatusResponse) GetFirstEvaluationDateAsTime() (*time.Time, error) {
	if o.FirstEvaluationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FirstEvaluationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *AssessmentStatusResponse) SetFirstEvaluationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FirstEvaluationDate = &formatted
}

func (o *AssessmentStatusResponse) GetStatusChangeDateAsTime() (*time.Time, error) {
	if o.StatusChangeDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StatusChangeDate, "2006-01-02T15:04:05Z07:00")
}

func (o *AssessmentStatusResponse) SetStatusChangeDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StatusChangeDate = &formatted
}
