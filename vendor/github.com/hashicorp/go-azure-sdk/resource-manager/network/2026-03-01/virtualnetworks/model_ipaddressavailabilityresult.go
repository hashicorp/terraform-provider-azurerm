package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPAddressAvailabilityResult struct {
	Available            *bool     `json:"available,omitempty"`
	AvailableIPAddresses *[]string `json:"availableIPAddresses,omitempty"`
	IsPlatformReserved   *bool     `json:"isPlatformReserved,omitempty"`
}
