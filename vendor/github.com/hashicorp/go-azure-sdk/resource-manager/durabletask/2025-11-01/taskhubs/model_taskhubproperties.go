package taskhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskHubProperties struct {
	DashboardURL      *string            `json:"dashboardUrl,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
