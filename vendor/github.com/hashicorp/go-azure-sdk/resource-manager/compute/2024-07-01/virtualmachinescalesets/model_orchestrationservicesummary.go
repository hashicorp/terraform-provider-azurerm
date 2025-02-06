package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrchestrationServiceSummary struct {
	ServiceName  *OrchestrationServiceNames `json:"serviceName,omitempty"`
	ServiceState *OrchestrationServiceState `json:"serviceState,omitempty"`
}
