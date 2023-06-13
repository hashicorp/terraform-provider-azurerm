package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPv6Config struct {
	Gateway      *string `json:"gateway,omitempty"`
	IPAddress    *string `json:"ipAddress,omitempty"`
	PrefixLength *int64  `json:"prefixLength,omitempty"`
}
