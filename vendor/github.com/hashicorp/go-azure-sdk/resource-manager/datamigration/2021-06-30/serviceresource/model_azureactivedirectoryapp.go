package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureActiveDirectoryApp struct {
	AppKey        string `json:"appKey"`
	ApplicationId string `json:"applicationId"`
	TenantId      string `json:"tenantId"`
}
