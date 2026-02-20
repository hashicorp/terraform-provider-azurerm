package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterface struct {
	IPAddresses *[]IPAddress `json:"ipAddresses,omitempty"`
	Id          *string      `json:"id,omitempty"`
	MacAddress  *string      `json:"macAddress,omitempty"`
	Name        *string      `json:"name,omitempty"`
}
