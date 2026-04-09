package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2ErrorSettings struct {
	ResendIfMDNNotReceived  bool `json:"resendIfMDNNotReceived"`
	SuspendDuplicateMessage bool `json:"suspendDuplicateMessage"`
}
