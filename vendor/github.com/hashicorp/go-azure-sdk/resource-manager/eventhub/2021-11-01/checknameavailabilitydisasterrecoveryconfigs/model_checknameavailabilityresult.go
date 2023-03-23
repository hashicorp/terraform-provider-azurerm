package checknameavailabilitydisasterrecoveryconfigs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityResult struct {
	Message       *string            `json:"message,omitempty"`
	NameAvailable *bool              `json:"nameAvailable,omitempty"`
	Reason        *UnavailableReason `json:"reason,omitempty"`
}
