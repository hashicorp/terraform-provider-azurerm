package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceConnectivityEndpoints struct {
	PrivateIPAddress *string `json:"privateIpAddress,omitempty"`
	PublicIPAddress  *string `json:"publicIpAddress,omitempty"`
}
