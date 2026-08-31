package networkmanagerconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerConnectionProperties struct {
	ConnectionState  *ScopeConnectionState `json:"connectionState,omitempty"`
	Description      *string               `json:"description,omitempty"`
	NetworkManagerId *string               `json:"networkManagerId,omitempty"`
}
