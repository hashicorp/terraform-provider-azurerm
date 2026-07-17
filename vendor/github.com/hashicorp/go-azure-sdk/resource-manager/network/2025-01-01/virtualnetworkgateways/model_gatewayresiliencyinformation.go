package virtualnetworkgateways

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayResiliencyInformation struct {
	Components                  *[]ResiliencyRecommendationComponents `json:"components,omitempty"`
	LastComputedTime            *string                               `json:"lastComputedTime,omitempty"`
	MaxScoreFromRecommendations *string                               `json:"maxScoreFromRecommendations,omitempty"`
	MinScoreFromRecommendations *string                               `json:"minScoreFromRecommendations,omitempty"`
	NextEligibleComputeTime     *string                               `json:"nextEligibleComputeTime,omitempty"`
	OverallScore                *string                               `json:"overallScore,omitempty"`
	ScoreChange                 *string                               `json:"scoreChange,omitempty"`
}

func (o *GatewayResiliencyInformation) GetLastComputedTimeAsTime() (*time.Time, error) {
	if o.LastComputedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastComputedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GatewayResiliencyInformation) SetLastComputedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastComputedTime = &formatted
}

func (o *GatewayResiliencyInformation) GetNextEligibleComputeTimeAsTime() (*time.Time, error) {
	if o.NextEligibleComputeTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextEligibleComputeTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GatewayResiliencyInformation) SetNextEligibleComputeTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextEligibleComputeTime = &formatted
}
