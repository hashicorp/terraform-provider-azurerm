package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrchestrationServiceStateInput struct {
	Action      OrchestrationServiceStateAction `json:"action"`
	ServiceName OrchestrationServiceNames       `json:"serviceName"`
}
