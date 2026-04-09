package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemService struct {
	PublicIPAddress   *string `json:"publicIpAddress,omitempty"`
	SystemServiceType *string `json:"systemServiceType,omitempty"`
	Version           *string `json:"version,omitempty"`
}
