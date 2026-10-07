package reachabilityanalysisruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReachabilityAnalysisRunProperties struct {
	AnalysisResult    *string            `json:"analysisResult,omitempty"`
	Description       *string            `json:"description,omitempty"`
	ErrorMessage      *string            `json:"errorMessage,omitempty"`
	IntentContent     *IntentContent     `json:"intentContent,omitempty"`
	IntentId          string             `json:"intentId"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
