package portalconfig

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalConfigProperties struct {
	Cors            *PortalConfigCorsProperties       `json:"cors,omitempty"`
	Csp             *PortalConfigCspProperties        `json:"csp,omitempty"`
	Delegation      *PortalConfigDelegationProperties `json:"delegation,omitempty"`
	EnableBasicAuth *bool                             `json:"enableBasicAuth,omitempty"`
	Signin          *PortalConfigPropertiesSignin     `json:"signin,omitempty"`
	Signup          *PortalConfigPropertiesSignup     `json:"signup,omitempty"`
}
