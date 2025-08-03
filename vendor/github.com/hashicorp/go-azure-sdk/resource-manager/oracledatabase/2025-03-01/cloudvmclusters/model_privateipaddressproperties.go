package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateIPAddressProperties struct {
	DisplayName   string `json:"displayName"`
	HostnameLabel string `json:"hostnameLabel"`
	IPAddress     string `json:"ipAddress"`
	Ocid          string `json:"ocid"`
	SubnetId      string `json:"subnetId"`
}
