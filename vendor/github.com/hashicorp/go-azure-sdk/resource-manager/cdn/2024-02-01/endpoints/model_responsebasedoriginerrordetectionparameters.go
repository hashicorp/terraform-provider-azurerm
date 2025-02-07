package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResponseBasedOriginErrorDetectionParameters struct {
	HTTPErrorRanges                          *[]HTTPErrorRangeParameters      `json:"httpErrorRanges,omitempty"`
	ResponseBasedDetectedErrorTypes          *ResponseBasedDetectedErrorTypes `json:"responseBasedDetectedErrorTypes,omitempty"`
	ResponseBasedFailoverThresholdPercentage *int64                           `json:"responseBasedFailoverThresholdPercentage,omitempty"`
}
