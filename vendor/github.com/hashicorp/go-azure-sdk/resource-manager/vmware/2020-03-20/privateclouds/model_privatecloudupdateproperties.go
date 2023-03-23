package privateclouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateCloudUpdateProperties struct {
	IdentitySources   *[]IdentitySource  `json:"identitySources,omitempty"`
	Internet          *InternetEnum      `json:"internet,omitempty"`
	ManagementCluster *ManagementCluster `json:"managementCluster,omitempty"`
}
