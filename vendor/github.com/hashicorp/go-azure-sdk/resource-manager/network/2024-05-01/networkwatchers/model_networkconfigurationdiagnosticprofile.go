package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConfigurationDiagnosticProfile struct {
	Destination     string    `json:"destination"`
	DestinationPort string    `json:"destinationPort"`
	Direction       Direction `json:"direction"`
	Protocol        string    `json:"protocol"`
	Source          string    `json:"source"`
}
