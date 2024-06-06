package nginxcertificate

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxCertificateProperties struct {
	CertificateError       *NginxCertificateErrorResponseBody `json:"certificateError,omitempty"`
	CertificateVirtualPath *string                            `json:"certificateVirtualPath,omitempty"`
	KeyVaultSecretCreated  *string                            `json:"keyVaultSecretCreated,omitempty"`
	KeyVaultSecretId       *string                            `json:"keyVaultSecretId,omitempty"`
	KeyVaultSecretVersion  *string                            `json:"keyVaultSecretVersion,omitempty"`
	KeyVirtualPath         *string                            `json:"keyVirtualPath,omitempty"`
	ProvisioningState      *ProvisioningState                 `json:"provisioningState,omitempty"`
	Sha1Thumbprint         *string                            `json:"sha1Thumbprint,omitempty"`
}

func (o *NginxCertificateProperties) GetKeyVaultSecretCreatedAsTime() (*time.Time, error) {
	if o.KeyVaultSecretCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.KeyVaultSecretCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *NginxCertificateProperties) SetKeyVaultSecretCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.KeyVaultSecretCreated = &formatted
}
