package cloudhsmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudHsmProperties struct {
	Fqdn         *string `json:"fqdn,omitempty"`
	State        *string `json:"state,omitempty"`
	StateMessage *string `json:"stateMessage,omitempty"`
}
