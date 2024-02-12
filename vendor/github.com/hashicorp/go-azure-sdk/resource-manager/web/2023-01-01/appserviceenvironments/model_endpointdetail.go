package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointDetail struct {
	IPAddress    *string  `json:"ipAddress,omitempty"`
	IsAccessible *bool    `json:"isAccessible,omitempty"`
	Latency      *float64 `json:"latency,omitempty"`
	Port         *int64   `json:"port,omitempty"`
}
