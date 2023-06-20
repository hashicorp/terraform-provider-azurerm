package tenants

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityRequest struct {
	CountryCode *string `json:"countryCode,omitempty"`
	Name        *string `json:"name,omitempty"`
}
