package integrationruntimedisableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeVNetProperties struct {
	PublicIPs *[]string `json:"publicIPs,omitempty"`
	Subnet    *string   `json:"subnet,omitempty"`
	SubnetId  *string   `json:"subnetId,omitempty"`
	VNetId    *string   `json:"vNetId,omitempty"`
}
