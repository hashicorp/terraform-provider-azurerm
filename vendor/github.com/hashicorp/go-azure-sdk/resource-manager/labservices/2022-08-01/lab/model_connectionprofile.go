package lab

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionProfile struct {
	ClientRdpAccess *ConnectionType `json:"clientRdpAccess,omitempty"`
	ClientSshAccess *ConnectionType `json:"clientSshAccess,omitempty"`
	WebRdpAccess    *ConnectionType `json:"webRdpAccess,omitempty"`
	WebSshAccess    *ConnectionType `json:"webSshAccess,omitempty"`
}
