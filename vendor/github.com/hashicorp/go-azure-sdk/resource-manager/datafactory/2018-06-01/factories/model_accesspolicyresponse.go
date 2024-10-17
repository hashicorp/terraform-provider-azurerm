package factories

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyResponse struct {
	AccessToken  *string           `json:"accessToken,omitempty"`
	DataPlaneURL *string           `json:"dataPlaneUrl,omitempty"`
	Policy       *UserAccessPolicy `json:"policy,omitempty"`
}
