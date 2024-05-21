package replicationnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Subnet struct {
	AddressList  *[]string `json:"addressList,omitempty"`
	FriendlyName *string   `json:"friendlyName,omitempty"`
	Name         *string   `json:"name,omitempty"`
}
