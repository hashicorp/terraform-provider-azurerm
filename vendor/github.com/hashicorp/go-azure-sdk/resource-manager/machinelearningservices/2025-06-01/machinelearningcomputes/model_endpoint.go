package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Endpoint struct {
	HostIP    *string   `json:"hostIp,omitempty"`
	Name      *string   `json:"name,omitempty"`
	Protocol  *Protocol `json:"protocol,omitempty"`
	Published *int64    `json:"published,omitempty"`
	Target    *int64    `json:"target,omitempty"`
}
