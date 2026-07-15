package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModelCapacityCalculatorWorkloadRequestParam struct {
	AvgGeneratedTokens *int64 `json:"avgGeneratedTokens,omitempty"`
	AvgPromptTokens    *int64 `json:"avgPromptTokens,omitempty"`
}
