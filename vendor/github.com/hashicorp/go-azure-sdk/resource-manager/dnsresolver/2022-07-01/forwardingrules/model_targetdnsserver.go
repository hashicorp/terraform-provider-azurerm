package forwardingrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetDnsServer struct {
	IPAddress string `json:"ipAddress"`
	Port      *int64 `json:"port,omitempty"`
}
