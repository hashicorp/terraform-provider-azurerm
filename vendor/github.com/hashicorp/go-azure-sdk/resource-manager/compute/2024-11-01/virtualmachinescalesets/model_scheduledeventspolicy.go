package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledEventsPolicy struct {
	ScheduledEventsAdditionalPublishingTargets *ScheduledEventsAdditionalPublishingTargets `json:"scheduledEventsAdditionalPublishingTargets,omitempty"`
	UserInitiatedReboot                        *UserInitiatedReboot                        `json:"userInitiatedReboot,omitempty"`
	UserInitiatedRedeploy                      *UserInitiatedRedeploy                      `json:"userInitiatedRedeploy,omitempty"`
}
