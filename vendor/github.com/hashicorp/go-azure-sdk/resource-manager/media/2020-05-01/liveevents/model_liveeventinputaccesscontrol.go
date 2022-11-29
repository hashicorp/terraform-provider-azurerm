package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventInputAccessControl struct {
	IP *IPAccessControl `json:"ip,omitempty"`
}
