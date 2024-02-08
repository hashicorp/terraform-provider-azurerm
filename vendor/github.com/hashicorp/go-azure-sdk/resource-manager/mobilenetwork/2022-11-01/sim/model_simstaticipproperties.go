package sim

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SimStaticIPProperties struct {
	AttachedDataNetwork *AttachedDataNetworkResourceId `json:"attachedDataNetwork,omitempty"`
	Slice               *SliceResourceId               `json:"slice,omitempty"`
	StaticIP            *SimStaticIPPropertiesStaticIP `json:"staticIp,omitempty"`
}
