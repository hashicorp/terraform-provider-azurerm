package environments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngressEnvironmentStatus struct {
	State        *IngressState            `json:"state,omitempty"`
	StateDetails *EnvironmentStateDetails `json:"stateDetails,omitempty"`
}
