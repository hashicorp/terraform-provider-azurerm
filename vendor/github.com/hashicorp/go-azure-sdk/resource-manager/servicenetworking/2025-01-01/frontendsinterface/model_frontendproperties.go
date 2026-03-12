package frontendsinterface

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FrontendProperties struct {
	Fqdn              *string            `json:"fqdn,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
