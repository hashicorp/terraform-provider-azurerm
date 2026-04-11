package monitoredsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionList struct {
	MonitoredSubscriptionList *[]MonitoredSubscription `json:"monitoredSubscriptionList,omitempty"`
	PatchOperation            *PatchOperation          `json:"patchOperation,omitempty"`
	ProvisioningState         *ProvisioningState       `json:"provisioningState,omitempty"`
}
