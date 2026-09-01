package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngressPortMapping struct {
	ExposedPort *int64 `json:"exposedPort,omitempty"`
	External    bool   `json:"external"`
	TargetPort  int64  `json:"targetPort"`
}
