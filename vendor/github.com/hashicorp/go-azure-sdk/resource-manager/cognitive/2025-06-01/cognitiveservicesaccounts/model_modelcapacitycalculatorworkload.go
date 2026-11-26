package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModelCapacityCalculatorWorkload struct {
	RequestParameters *ModelCapacityCalculatorWorkloadRequestParam `json:"requestParameters,omitempty"`
	RequestPerMinute  *int64                                       `json:"requestPerMinute,omitempty"`
}
