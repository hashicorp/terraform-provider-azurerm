package replicationnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProperties struct {
	FabricType   *string   `json:"fabricType,omitempty"`
	FriendlyName *string   `json:"friendlyName,omitempty"`
	NetworkType  *string   `json:"networkType,omitempty"`
	Subnets      *[]Subnet `json:"subnets,omitempty"`
}
