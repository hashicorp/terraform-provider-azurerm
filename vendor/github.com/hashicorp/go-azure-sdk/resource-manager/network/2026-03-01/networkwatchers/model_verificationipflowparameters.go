package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerificationIPFlowParameters struct {
	Direction           Direction      `json:"direction"`
	LocalIPAddress      string         `json:"localIPAddress"`
	LocalPort           string         `json:"localPort"`
	Protocol            IPFlowProtocol `json:"protocol"`
	RemoteIPAddress     string         `json:"remoteIPAddress"`
	RemotePort          string         `json:"remotePort"`
	TargetNicResourceId *string        `json:"targetNicResourceId,omitempty"`
	TargetResourceId    string         `json:"targetResourceId"`
}
