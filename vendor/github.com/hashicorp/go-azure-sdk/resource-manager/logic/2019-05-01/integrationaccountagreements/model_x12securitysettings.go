package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12SecuritySettings struct {
	AuthorizationQualifier string  `json:"authorizationQualifier"`
	AuthorizationValue     *string `json:"authorizationValue,omitempty"`
	PasswordValue          *string `json:"passwordValue,omitempty"`
	SecurityQualifier      string  `json:"securityQualifier"`
}
