package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkACL struct {
	Allow *[]SignalRRequestType `json:"allow,omitempty"`
	Deny  *[]SignalRRequestType `json:"deny,omitempty"`
}
