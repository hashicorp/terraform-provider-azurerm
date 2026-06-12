package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomain struct {
	BindingType   *BindingType `json:"bindingType,omitempty"`
	CertificateId *string      `json:"certificateId,omitempty"`
	Name          string       `json:"name"`
}
