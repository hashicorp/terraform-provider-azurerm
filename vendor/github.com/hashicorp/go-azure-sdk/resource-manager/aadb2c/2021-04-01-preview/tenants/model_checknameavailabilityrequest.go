package tenants

// Copyright IBM Corp. 2021, 2025 All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityRequest struct {
	CountryCode *string `json:"countryCode,omitempty"`
	Name        *string `json:"name,omitempty"`
}
