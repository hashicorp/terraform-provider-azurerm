package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ControlPlaneProfile struct {
	ControlPlaneEndpoint *ControlPlaneProfileControlPlaneEndpoint `json:"controlPlaneEndpoint,omitempty"`
	Count                *int64                                   `json:"count,omitempty"`
	VMSize               *string                                  `json:"vmSize,omitempty"`
}
