package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPropertiesVMipPoolInlined struct {
	EndIP   *string `json:"endIP,omitempty"`
	StartIP *string `json:"startIP,omitempty"`
}
