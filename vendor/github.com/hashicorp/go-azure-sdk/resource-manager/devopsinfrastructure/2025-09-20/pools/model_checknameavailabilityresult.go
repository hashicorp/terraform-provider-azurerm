package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityResult struct {
	Available AvailabilityStatus          `json:"available"`
	Message   string                      `json:"message"`
	Name      string                      `json:"name"`
	Reason    CheckNameAvailabilityReason `json:"reason"`
}
