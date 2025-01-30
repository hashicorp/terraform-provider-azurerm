package reachabilityanalysisintents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPTraffic struct {
	DestinationIPs   []string          `json:"destinationIps"`
	DestinationPorts []string          `json:"destinationPorts"`
	Protocols        []NetworkProtocol `json:"protocols"`
	SourceIPs        []string          `json:"sourceIps"`
	SourcePorts      []string          `json:"sourcePorts"`
}
