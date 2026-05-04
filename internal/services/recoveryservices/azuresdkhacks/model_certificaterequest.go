package azuresdkhacks

import (
	"encoding/json"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-04-01/vaultcertificates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateRequest struct {
	Properties    *vaultcertificates.RawCertificateData `json:"properties,omitempty"`
	CreateOptions *CertificateCreateOptions             `json:"certificateCreateOptions,omitempty"`
}

type CertificateCreateOptions struct {
	ValidityInHours int64 `json:"validityInHours,omitempty"`
}

func (c CertificateCreateOptions) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})

	objectMap["validityInHours"] = c.ValidityInHours
	return json.Marshal(objectMap)
}
