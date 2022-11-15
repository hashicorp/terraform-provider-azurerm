package dedicatedhsms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	NetworkInterfaces *[]NetworkInterface `json:"networkInterfaces,omitempty"`
	Subnet            *ApiEntityReference `json:"subnet,omitempty"`
}
