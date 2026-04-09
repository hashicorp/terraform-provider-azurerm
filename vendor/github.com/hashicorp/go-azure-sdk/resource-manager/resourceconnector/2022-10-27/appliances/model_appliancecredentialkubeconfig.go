package appliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplianceCredentialKubeconfig struct {
	Name  *AccessProfileType `json:"name,omitempty"`
	Value *string            `json:"value,omitempty"`
}
