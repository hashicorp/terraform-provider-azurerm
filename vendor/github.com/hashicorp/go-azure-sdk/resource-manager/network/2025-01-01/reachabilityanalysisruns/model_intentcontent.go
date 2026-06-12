package reachabilityanalysisruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntentContent struct {
	Description           *string   `json:"description,omitempty"`
	DestinationResourceId string    `json:"destinationResourceId"`
	IPTraffic             IPTraffic `json:"ipTraffic"`
	SourceResourceId      string    `json:"sourceResourceId"`
}
