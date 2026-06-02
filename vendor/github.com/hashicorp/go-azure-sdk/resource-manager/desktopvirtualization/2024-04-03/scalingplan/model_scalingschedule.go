package scalingplan

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScalingSchedule struct {
	DaysOfWeek                     *[]DaysOfWeek                      `json:"daysOfWeek,omitempty"`
	Name                           *string                            `json:"name,omitempty"`
	OffPeakLoadBalancingAlgorithm  *SessionHostLoadBalancingAlgorithm `json:"offPeakLoadBalancingAlgorithm,omitempty"`
	OffPeakStartTime               *Time                              `json:"offPeakStartTime,omitempty"`
	PeakLoadBalancingAlgorithm     *SessionHostLoadBalancingAlgorithm `json:"peakLoadBalancingAlgorithm,omitempty"`
	PeakStartTime                  *Time                              `json:"peakStartTime,omitempty"`
	RampDownCapacityThresholdPct   *int64                             `json:"rampDownCapacityThresholdPct,omitempty"`
	RampDownForceLogoffUsers       *bool                              `json:"rampDownForceLogoffUsers,omitempty"`
	RampDownLoadBalancingAlgorithm *SessionHostLoadBalancingAlgorithm `json:"rampDownLoadBalancingAlgorithm,omitempty"`
	RampDownMinimumHostsPct        *int64                             `json:"rampDownMinimumHostsPct,omitempty"`
	RampDownNotificationMessage    *string                            `json:"rampDownNotificationMessage,omitempty"`
	RampDownStartTime              *Time                              `json:"rampDownStartTime,omitempty"`
	RampDownStopHostsWhen          *StopHostsWhen                     `json:"rampDownStopHostsWhen,omitempty"`
	RampDownWaitTimeMinutes        *int64                             `json:"rampDownWaitTimeMinutes,omitempty"`
	RampUpCapacityThresholdPct     *int64                             `json:"rampUpCapacityThresholdPct,omitempty"`
	RampUpLoadBalancingAlgorithm   *SessionHostLoadBalancingAlgorithm `json:"rampUpLoadBalancingAlgorithm,omitempty"`
	RampUpMinimumHostsPct          *int64                             `json:"rampUpMinimumHostsPct,omitempty"`
	RampUpStartTime                *Time                              `json:"rampUpStartTime,omitempty"`
}
