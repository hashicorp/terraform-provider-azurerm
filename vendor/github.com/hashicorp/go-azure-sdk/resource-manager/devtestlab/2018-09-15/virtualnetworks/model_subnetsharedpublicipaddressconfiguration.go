package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubnetSharedPublicIPAddressConfiguration struct {
	AllowedPorts *[]Port `json:"allowedPorts,omitempty"`
}
