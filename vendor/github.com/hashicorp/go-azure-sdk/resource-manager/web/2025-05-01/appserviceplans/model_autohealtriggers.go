package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoHealTriggers struct {
	PrivateBytesInKB     *int64                          `json:"privateBytesInKB,omitempty"`
	Requests             *RequestsBasedTrigger           `json:"requests,omitempty"`
	SlowRequests         *SlowRequestsBasedTrigger       `json:"slowRequests,omitempty"`
	SlowRequestsWithPath *[]SlowRequestsBasedTrigger     `json:"slowRequestsWithPath,omitempty"`
	StatusCodes          *[]StatusCodesBasedTrigger      `json:"statusCodes,omitempty"`
	StatusCodesRange     *[]StatusCodesRangeBasedTrigger `json:"statusCodesRange,omitempty"`
}
