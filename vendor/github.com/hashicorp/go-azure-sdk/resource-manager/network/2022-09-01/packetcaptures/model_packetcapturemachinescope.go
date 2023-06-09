package packetcaptures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCaptureMachineScope struct {
	Exclude *[]string `json:"exclude,omitempty"`
	Include *[]string `json:"include,omitempty"`
}
