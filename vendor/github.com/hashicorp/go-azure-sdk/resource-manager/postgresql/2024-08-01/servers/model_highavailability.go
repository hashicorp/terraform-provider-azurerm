package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HighAvailability struct {
	Mode                    *HighAvailabilityMode `json:"mode,omitempty"`
	StandbyAvailabilityZone *string               `json:"standbyAvailabilityZone,omitempty"`
	State                   *ServerHAState        `json:"state,omitempty"`
}
