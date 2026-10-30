package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	IPAddresses          *[]string `json:"ipAddresses,omitempty"`
	StaticIPAddressCount *int64    `json:"staticIpAddressCount,omitempty"`
	SubnetId             *string   `json:"subnetId,omitempty"`
}
