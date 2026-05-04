package reachabilityanalysisintents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReachabilityAnalysisIntentProperties struct {
	Description           *string            `json:"description,omitempty"`
	DestinationResourceId string             `json:"destinationResourceId"`
	IPTraffic             IPTraffic          `json:"ipTraffic"`
	ProvisioningState     *ProvisioningState `json:"provisioningState,omitempty"`
	SourceResourceId      string             `json:"sourceResourceId"`
}
