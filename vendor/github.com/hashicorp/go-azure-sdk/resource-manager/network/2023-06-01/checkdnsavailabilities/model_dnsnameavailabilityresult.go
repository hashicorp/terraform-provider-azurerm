package checkdnsavailabilities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsNameAvailabilityResult struct {
	Available *bool `json:"available,omitempty"`
}
