package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NsgCidr struct {
	DestinationPortRange *PortRange `json:"destinationPortRange,omitempty"`
	Source               string     `json:"source"`
}
