package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityResponse struct {
	Available *bool                        `json:"available,omitempty"`
	Message   *string                      `json:"message,omitempty"`
	Name      *string                      `json:"name,omitempty"`
	Reason    *CheckNameAvailabilityReason `json:"reason,omitempty"`
}
