package availabilitygrouplisteners

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateIPAddress struct {
	IPAddress        *string `json:"ipAddress,omitempty"`
	SubnetResourceId *string `json:"subnetResourceId,omitempty"`
}
