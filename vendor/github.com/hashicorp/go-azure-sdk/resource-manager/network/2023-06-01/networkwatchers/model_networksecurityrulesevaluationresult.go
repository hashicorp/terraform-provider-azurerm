package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityRulesEvaluationResult struct {
	DestinationMatched     *bool   `json:"destinationMatched,omitempty"`
	DestinationPortMatched *bool   `json:"destinationPortMatched,omitempty"`
	Name                   *string `json:"name,omitempty"`
	ProtocolMatched        *bool   `json:"protocolMatched,omitempty"`
	SourceMatched          *bool   `json:"sourceMatched,omitempty"`
	SourcePortMatched      *bool   `json:"sourcePortMatched,omitempty"`
}
