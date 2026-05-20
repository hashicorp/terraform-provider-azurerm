package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JwtClaimChecks struct {
	AllowedClientApplications *[]string `json:"allowedClientApplications,omitempty"`
	AllowedGroups             *[]string `json:"allowedGroups,omitempty"`
}
