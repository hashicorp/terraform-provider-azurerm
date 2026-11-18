package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StatusCodesRangeBasedTrigger struct {
	Count        *int64  `json:"count,omitempty"`
	Path         *string `json:"path,omitempty"`
	StatusCodes  *string `json:"statusCodes,omitempty"`
	TimeInterval *string `json:"timeInterval,omitempty"`
}
