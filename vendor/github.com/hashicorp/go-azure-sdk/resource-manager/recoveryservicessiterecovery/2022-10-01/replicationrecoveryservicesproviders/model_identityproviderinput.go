package replicationrecoveryservicesproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityProviderInput struct {
	AadAuthority  string `json:"aadAuthority"`
	ApplicationId string `json:"applicationId"`
	Audience      string `json:"audience"`
	ObjectId      string `json:"objectId"`
	TenantId      string `json:"tenantId"`
}
