package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NextHopParameters struct {
	DestinationIPAddress string  `json:"destinationIPAddress"`
	SourceIPAddress      string  `json:"sourceIPAddress"`
	TargetNicResourceId  *string `json:"targetNicResourceId,omitempty"`
	TargetResourceId     string  `json:"targetResourceId"`
}
