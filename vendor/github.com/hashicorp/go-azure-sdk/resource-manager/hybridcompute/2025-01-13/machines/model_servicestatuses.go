package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceStatuses struct {
	ExtensionService          *ServiceStatus `json:"extensionService,omitempty"`
	GuestConfigurationService *ServiceStatus `json:"guestConfigurationService,omitempty"`
}
