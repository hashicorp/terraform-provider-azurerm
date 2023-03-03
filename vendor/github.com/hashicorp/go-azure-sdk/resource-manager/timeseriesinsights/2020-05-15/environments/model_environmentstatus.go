package environments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentStatus struct {
	Ingress     *IngressEnvironmentStatus     `json:"ingress,omitempty"`
	WarmStorage *WarmStorageEnvironmentStatus `json:"warmStorage,omitempty"`
}
